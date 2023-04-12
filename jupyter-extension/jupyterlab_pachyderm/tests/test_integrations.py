import os
import sys
import subprocess
import time
import json
from datetime import datetime
from pathlib import Path
from random import randint

import pytest
import requests

from jupyterlab_pachyderm.handlers import NAMESPACE, VERSION
from jupyterlab_pachyderm.env import PFS_MOUNT_DIR

from . import TEST_NOTEBOOK

ADDRESS = "http://localhost:8888"
BASE_URL = f"{ADDRESS}/{NAMESPACE}/{VERSION}"
CONFIG_PATH = "~/.pachyderm/config.json"
ROOT_TOKEN = "iamroot"
DEFAULT_PROJECT = "default"


@pytest.fixture()
def pachyderm_resources():
    print("creating pachyderm resources")
    import python_pachyderm
    from python_pachyderm.pfs import Commit

    repos = ["images", "edges", "montage"]
    branches = ["master", "dev"]
    files = ["file1", "file2"]

    client = python_pachyderm.Client()
    client.delete_all()

    for repo in repos:
        client.create_repo(repo)
        for branch in branches:
            for file in files:
                client.put_file_bytes(
                    Commit(repo=repo, branch=branch), file, value=b"some data"
                )

    yield repos, branches, files


@pytest.fixture(scope="module")
def dev_server(request):
    """If you are having trouble with this fixture,
    consider running the mount-server locally.

    Tests can set the marker mount_server to value False if they do not
      require the mount-server as a dependency.
    """
    marker = request.node.get_closest_marker("mount_server")
    use_mount_server = marker is None or marker.args[0]

    print("starting development server...")
    p = subprocess.Popen(
        [sys.executable, "-m", "jupyterlab_pachyderm.dev_server"],
        env={"PFS_MOUNT_DIR": PFS_MOUNT_DIR},
        stdout=subprocess.PIPE,
    )
    # Give time for python test server to start
    time.sleep(1)

    if use_mount_server:
        requests.put(
            f"{BASE_URL}/config",
            data=json.dumps({"pachd_address": "localhost:30650"})
        )

        # Give time for mount server to start
        for _ in range(15):
            try:
                r = requests.get(f"{BASE_URL}/config", timeout=1)
                if r.status_code == 200 and r.json()["cluster_status"] != "INVALID":
                    break
            except Exception:
                pass
            time.sleep(1)
        else:  # If the for-loop isn't escaped.
            subprocess.run(["pkill", "-f", "mount-server"])
            raise RuntimeError("mount server is having issues starting up")

    yield

    print("killing development server...")
    subprocess.run(["pkill", "-f", "mount-server"])
    subprocess.run(
        ["bash", "-c", f"umount {PFS_MOUNT_DIR}"],
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
    )
    p.terminate()
    p.wait()


def test_list_repos(pachyderm_resources, dev_server):
    repos, branches, _ = pachyderm_resources

    r = requests.get(f"{BASE_URL}/repos")

    assert r.status_code == 200
    for repo_info in r.json():
        assert repo_info.keys() == {"authorization", "branches", "repo", "project"}
        assert repo_info["repo"] in repos
        for _branch in repo_info["branches"]:
            assert _branch in branches


def test_list_mounts(pachyderm_resources, dev_server):
    repos, branches, _ = pachyderm_resources

    r = requests.put(
        f"{BASE_URL}/_mount",
        data=json.dumps(
            {
                "mounts": [
                    {
                        "name": "mount1",
                        "repo": repos[0],
                    }
                ]
            }
        ),
    )
    assert r.status_code == 200

    r = requests.get(f"{BASE_URL}/mounts")
    assert r.status_code == 200

    resp = r.json()
    assert len(resp["mounted"]) == 1
    for mount_info in resp["mounted"]:
        assert mount_info.keys() == {
            "name",
            "repo",
            "branch",
            "project",
            "commit",
            "files",
            "mode",
            "state",
            "status",
            "mountpoint",
            "actual_mounted_commit",
            "latest_commit",
            "how_many_commits_behind",
        }

    for _repo_info in resp["unmounted"]:
        assert _repo_info["repo"] in repos
        assert _repo_info.keys() == {"authorization", "branches", "repo", "project"}
        for _branch in _repo_info["branches"]:
            assert _branch in branches
    assert len(resp["unmounted"]) == len(repos)


def test_mount(pachyderm_resources, dev_server):
    repos, _, files = pachyderm_resources

    to_mount = {
        "mounts": [
            {
                "name": repos[0],
                "repo": repos[0],
                "branch": "master",
                "project": DEFAULT_PROJECT,
                "mode": "ro",
            },
            {
                "name": repos[0] + "_dev",
                "repo": repos[0],
                "branch": "dev",
                "project": DEFAULT_PROJECT,
                "mode": "ro",
            },
            {
                "name": repos[1],
                "repo": repos[1],
                "branch": "master",
                "project": DEFAULT_PROJECT,
                "mode": "ro",
            },
        ]
    }
    r = requests.put(f"{BASE_URL}/_mount", data=json.dumps(to_mount))
    assert r.status_code == 200

    resp = r.json()
    assert len(resp["mounted"]) == 3
    assert len(list(os.walk(PFS_MOUNT_DIR))[0][1]) == 3
    for mount_info in resp["mounted"]:
        assert sorted(
            list(os.walk(os.path.join(PFS_MOUNT_DIR, mount_info["name"])))[0][2]
        ) == sorted(files)
    assert len(resp["unmounted"]) == 3
    assert len(resp["unmounted"][0]["branches"]) == 2
    assert len(resp["unmounted"][1]["branches"]) == 2

    r = requests.put(
        f"{BASE_URL}/_unmount_all",
    )
    assert r.status_code == 200
    assert r.json()["mounted"] == []
    assert len(r.json()["unmounted"]) == 3
    assert list(os.walk(PFS_MOUNT_DIR)) == [(PFS_MOUNT_DIR, [], [])]


def test_unmount(pachyderm_resources, dev_server):
    repos, branches, files = pachyderm_resources

    to_mount = {
        "mounts": [
            {
                "name": repos[0],
                "repo": repos[0],
                "branch": "master",
                "project": DEFAULT_PROJECT,
                "mode": "ro",
            },
            {
                "name": repos[0] + "_dev",
                "repo": repos[0],
                "branch": "dev",
                "project": DEFAULT_PROJECT,
                "mode": "ro",
            },
        ]
    }
    r = requests.put(f"{BASE_URL}/_mount", data=json.dumps(to_mount))
    assert r.status_code == 200
    assert sorted(list(os.walk(os.path.join(PFS_MOUNT_DIR, repos[0])))[0][2]) == sorted(
        files
    )
    assert sorted(
        list(os.walk(os.path.join(PFS_MOUNT_DIR, repos[0] + "_dev")))[0][2]
    ) == sorted(files)
    assert len(r.json()["mounted"]) == 2
    assert len(r.json()["unmounted"]) == 3

    r = requests.put(
        f"{BASE_URL}/_unmount",
        data=json.dumps({"mounts": [repos[0] + "_dev"]}),
    )
    assert r.status_code == 200
    assert len(r.json()["mounted"]) == 1
    assert len(r.json()["unmounted"]) == 3
    assert len(r.json()["unmounted"][0]["branches"]) == 2

    r = requests.put(
        f"{BASE_URL}/_unmount",
        data=json.dumps({"mounts": [repos[0]]}),
    )
    assert r.status_code == 200
    assert len(r.json()["mounted"]) == 0
    assert len(r.json()["unmounted"]) == 3
    assert len(r.json()["unmounted"][0]["branches"]) == 2
    assert list(os.walk(PFS_MOUNT_DIR)) == [(PFS_MOUNT_DIR, [], [])]


def test_mount_datums(pachyderm_resources, dev_server):
    repos, branches, files = pachyderm_resources
    input_spec = {
        "input": {
            "cross": [
                {
                    "pfs": {
                        "repo": repos[0],
                        "glob": "/",
                    }
                },
                {
                    "pfs": {
                        "repo": repos[1],
                        "branch": "dev",
                        "glob": "/*",
                    }
                },
                {
                    "pfs": {
                        "repo": repos[2],
                        "glob": "/*",
                    }
                },
            ]
        }
    }

    r = requests.put(f"{BASE_URL}/_mount_datums", data=json.dumps(input_spec))
    assert r.status_code == 200
    assert r.json()["idx"] == 0
    assert r.json()["num_datums"] == 4
    list(os.walk(os.path.join(PFS_MOUNT_DIR, "out")))  # makes "out" dir appear
    assert len(list(os.walk(PFS_MOUNT_DIR))[0][1]) == 4
    datum0_id = r.json()["id"]

    assert sorted(list(os.walk(os.path.join(PFS_MOUNT_DIR, "".join([DEFAULT_PROJECT, "_", repos[0]]))))[0][2]) == sorted(
        files
    )
    assert "".join([DEFAULT_PROJECT, "_", repos[1], "_dev"]) in list(os.walk(PFS_MOUNT_DIR))[0][1]
    assert len(list(os.walk(os.path.join(PFS_MOUNT_DIR, "".join([DEFAULT_PROJECT, "_", repos[2]]))))[0][2]) == 1

    r = requests.put(f"{BASE_URL}/_show_datum", params={"idx": "2"})
    assert r.status_code == 200
    assert r.json()["idx"] == 2
    assert r.json()["num_datums"] == 4
    assert r.json()["id"] != datum0_id

    r = requests.put(
        f"{BASE_URL}/_show_datum",
        params={
            "idx": "2",
            "id": datum0_id,
        },
    )
    assert r.status_code == 200
    assert r.json()["idx"] == 0
    assert r.json()["num_datums"] == 4
    assert r.json()["id"] == datum0_id

    r = requests.get(f"{BASE_URL}/datums")
    assert r.status_code == 200
    assert r.json()["input"] == {
        "cross": [
            {
                "pfs": {
                    "repo": repos[0],
                    "name": "".join([DEFAULT_PROJECT, "_", repos[0]]),
                    "glob": "/",
                    "project": DEFAULT_PROJECT,
                    "branch": "master",
                }
            },
            {
                "pfs": {
                    "repo": repos[1],
                    "name": "".join([DEFAULT_PROJECT, "_", repos[1], "_dev"]),
                    "branch": "dev",
                    "glob": "/*",
                    "project": DEFAULT_PROJECT,
                }
            },
            {
                "pfs": {
                    "repo": repos[2],
                    "name": "".join([DEFAULT_PROJECT, "_", repos[2]]),
                    "glob": "/*",
                    "project": DEFAULT_PROJECT,
                    "branch": "master",
                }
            },
        ]
    }
    assert r.json()["num_datums"] == 4
    assert r.json()["curr_idx"] == 0

    r = requests.put(f"{BASE_URL}/_unmount_all")
    assert r.status_code == 200


def test_config(dev_server):
    # PUT request
    test_endpoint = "localhost:30650"
    r = requests.put(
        f"{BASE_URL}/config", data=json.dumps({"pachd_address": test_endpoint})
    )

    config = json.load(open(os.path.expanduser(CONFIG_PATH)))
    active_context = config["v2"]["active_context"]
    try:
        endpoint_in_config = config["v2"]["contexts"][active_context]["pachd_address"]
    except:
        endpoint_in_config = str(
            config["v2"]["contexts"][active_context]["port_forwarders"]["pachd"]
        )

    assert r.status_code == 200
    assert r.json()["cluster_status"] != "INVALID"
    assert "30650" in endpoint_in_config

    # GET request
    r = requests.get(f"{BASE_URL}/config")

    assert r.status_code == 200
    assert r.json()["cluster_status"] != "INVALID"


@pytest.fixture
def simple_pachyderm_env():
    from python_pachyderm import Client
    client = Client()

    suffix = str(randint(100000, 999999))
    repo_name = f"images_{suffix}"
    pipeline_name = f"test_pipeline_{suffix}"
    client.delete_repo(repo_name, force=True)
    client.create_repo(repo_name)
    yield client, repo_name, pipeline_name
    client.delete_pipeline(pipeline_name, force=True)
    client.delete_pipeline(f"{pipeline_name}__context", force=True)
    client.delete_repo(repo_name, force=True)


@pytest.fixture
def notebook_path(simple_pachyderm_env) -> Path:
    """Yields a path to a notebook file suitable for testing.

    This writes a temporary notebook file with its metadata populated
      with the expected pipeline and repo names provided by the
      simple_pachyderm_env fixture.
    """
    _client, repo_name, pipeline_name = simple_pachyderm_env

    # Do a considerable amount of data munging.
    notebook_data = json.loads(TEST_NOTEBOOK.read_bytes())
    notebook_data['metadata']['same_config']['metadata']['name'] = pipeline_name
    input_spec = json.loads(notebook_data['metadata']['same_config']['run']['input'])
    input_spec['pfs']['repo'] = repo_name
    notebook_data['metadata']['same_config']['run']['input'] = json.dumps(input_spec)

    notebook_path = TEST_NOTEBOOK.with_stem(f"{TEST_NOTEBOOK.stem}_generated")
    notebook_path.write_text(json.dumps(notebook_data))
    yield notebook_path.relative_to(os.getcwd())
    if notebook_path.exists():
        notebook_path.unlink()


@pytest.mark.mount_server(False)
@pytest.mark.usefixtures('dev_server')
class TestPpsClient:

    @staticmethod
    def test_pps(simple_pachyderm_env, notebook_path):
        client, repo_name, pipeline_name = simple_pachyderm_env
        last_modified = datetime.utcfromtimestamp(os.path.getmtime(notebook_path))
        data = dict(last_modified_time=f"{datetime.isoformat(last_modified)}Z")
        r = requests.put(f"{BASE_URL}/pps/_create/{notebook_path}", data=json.dumps(data))
        assert r.status_code == 200
        assert next(client.inspect_pipeline(pipeline_name))
        assert r.json()["message"] == ("Create pipeline request sent. You may monitor its "
        "status by running \"pachctl list pipelines\" in a terminal.")

    @staticmethod
    def test_pps_validation_errors(notebook_path):
        r = requests.put(f"{BASE_URL}/pps/_create/{notebook_path}", data=json.dumps({}))
        assert r.status_code == 500
        assert r.json()['reason'] == f"Bad Request: last_modified_time not specified"
