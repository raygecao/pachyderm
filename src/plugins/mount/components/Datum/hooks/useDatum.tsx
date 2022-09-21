import {requestAPI} from '../../../../../handler';
import {useEffect, useState} from 'react';
import {ServerConnection} from '@jupyterlab/services';
import {DatumsResponse} from 'plugins/mount/types';

export type useDatumResponse = {
  loading: boolean;
  shouldShowCycler: boolean;
  currentDatumId: string;
  currentDatumIdx: number;
  setCurrentDatumIdx: (idx: number) => void;
  numDatums: number;
  inputSpec: string;
  setInputSpec: (input: string) => void;
  callMountDatums: () => Promise<void>;
  callUnmountAll: () => Promise<void>;
  errorMessage: string;
  setDebug: (text: string) => void;
};

export const useDatum = (
  showDatum: boolean,
  keepMounted: boolean,
  refresh: () => void,
  pollRefresh: () => Promise<void>,
  currentDatumInfo?: DatumsResponse,
): useDatumResponse => {
  const [loading, setLoading] = useState(false);
  const [shouldShowCycler, setShouldShowCycler] = useState(false);
  const [currentDatumId, setCurrentDatumId] = useState('');
  const [currentDatumIdx, setCurrentDatumIdx] = useState(-1);
  const [numDatums, setNumDatums] = useState(-1);
  const [inputSpec, setInputSpec] = useState('');
  const [errorMessage, setErrorMessage] = useState('');
  const [debug, setDebug] = useState('');

  useEffect(() => {
    if (showDatum && currentDatumIdx !== -1) {
      callShowDatum();
    }
  }, [currentDatumIdx, showDatum]);

  useEffect(() => {
    if (showDatum && !keepMounted) {
      callUnmountAll();
    }
    if (keepMounted && currentDatumInfo) {
      setShouldShowCycler(true);
      setCurrentDatumIdx(currentDatumInfo.curr_idx);
      setNumDatums(currentDatumInfo.num_datums);
      setInputSpec(JSON.stringify(currentDatumInfo.input, null, 2));
    }
  }, [showDatum]);

  const callMountDatums = async () => {
    setLoading(true);
    setErrorMessage('');
    console.log('==========>> DeBUG ', debug);
    console.log('==========>>> YOOOOO ', inputSpec);
    console.log('==========>>> len is ', inputSpec.length);

    try {
      const res = await requestAPI<any>('_mount_datums', 'PUT', {
        input: JSON.parse(inputSpec),
      });
      refresh();
      setCurrentDatumId(res.id);
      setCurrentDatumIdx(res.idx);
      setNumDatums(res.num_datums);
      setShouldShowCycler(true);
      setInputSpec(JSON.stringify(JSON.parse(inputSpec), null, 2));
    } catch (e) {
      console.log(e);
      if (e instanceof ServerConnection.ResponseError) {
        setErrorMessage('Bad data in input spec');
      } else if (e instanceof SyntaxError) {
        setErrorMessage('Poorly formatted json input spec');
      } else {
        setErrorMessage('Error mounting datums');
      }
    }

    setLoading(false);
  };

  const callShowDatum = async () => {
    setLoading(true);

    try {
      const res = await requestAPI<any>(
        `_show_datum?idx=${currentDatumIdx}`,
        'PUT',
      );
      refresh();
      setCurrentDatumId(res.id);
    } catch (e) {
      console.log(e);
    }

    setLoading(false);
  };

  const callUnmountAll = async () => {
    setLoading(true);

    try {
      refresh();
      await requestAPI<any>('_unmount_all', 'PUT');
      refresh();
      await pollRefresh();
      setCurrentDatumId('');
      setCurrentDatumIdx(-1);
      setNumDatums(0);
      setShouldShowCycler(false);
    } catch (e) {
      console.log(e);
    }

    setLoading(false);
  };

  return {
    loading,
    shouldShowCycler,
    currentDatumId,
    currentDatumIdx,
    setCurrentDatumIdx,
    numDatums,
    inputSpec,
    setInputSpec,
    callMountDatums,
    callUnmountAll,
    errorMessage,
    setDebug,
  };
};
