import * as React from 'react';

const SvgGenericError = (props) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width={400}
    height={180}
    viewBox="320 -20 180 270"
    xmlnsXlink="http://www.w3.org/1999/xlink"
    {...props}
  >
    <defs>
      <path id="genericError_svg__b" d="m197.12 159.04 29.68 53.2h-59.36z" />
      <path
        d="M10.328 21.474S-2.853 8.286 65.634.379C81.999-1.51 90.868 3.838 112.578 12.5c36.95 14.745 24.243 55.33 6.414 95.06-11.501 25.63-4.968 28.555-6.252 32.503-1.283 3.948-8.707 3.948-15.164 3.948-6.456 0-6.76-2.39-6.364-3.948.246-.976-1.33-13.27.812-29.328.793-5.951.793-9.918 2.776-16.663 0 0-8.443 1.83-13.484 1.587-14.277-.69-8.17 49.34-8.546 51.573-.793 4.682-7.515 5.59-13.068 4.967-5.552-.623-11.411-.04-10.221-3.158 1.88-4.93 4.868-54.175-5.047-54.175-18.243 0-.504 52.136-1.376 54.175-1.19 2.777-8.104 2.952-15.183 2.952-4.787 0-7.212-.397-6.816-4.76.179-1.984.77-8.72-2.752-16.857-5.646-13.045-5.013-28.873-5.6-36.304-.98-12.5-26.48-44.305-2.38-72.598z"
        id="genericError_svg__c"
      />
      <path
        d="M.233 11.923c-.67-2.18 2.637-5.33 4.313-6.876C7.318 2.492 15.874 2.03 21.823.84c17.35-3.471 30.886 8.239 35.072 23.864C60.3 37.418 52.657 38.986 44.687 51.65s-.477 20.622-7.089 24.448c-5.022 2.906-9.033-2.074-15.775-3.661-5.794-1.365-4.648-11.303-7.932-9.918-4.445 1.868-4.414-19.836-3.569-24.597l-10.089-26z"
        id="genericError_svg__e"
      />
      <path
        d="M58.575 14.09c.9-2.021.648-4.37-.659-6.154C55.76 5.072 50.524 2.263 45.13.898c-9.92-2.506-25.258 1.039-32.977 3.595C2.826 7.582.052 10.526.099 19.227c.052 9.364 8.586 16.238 8.784 24.271.189 7.612.734 11.26 6.293 23.491 5.559 12.232 4.641-4.74 18.818-5.044 5.451-.117 2.503-1.936 5.29-.167 3.768 2.38 6.715-18.832 6.6-23.576l12.69-24.113z"
        id="genericError_svg__g"
      />
      <path
        d="M83.82 79.64c-6.345 2.217-19.88 15.57-24.786 18.843-9.518 6.347-8.296 29.201-26.89 26.982-19.216-2.292-26.467-17.343-29.19-24.627-.922-2.466-5.092-12.956.012-10.669 5.103 2.287 8.686-7.897 9.395-1.504 1.477 13.317 10.228 23.042 17.365 22.005 7.638-1.11 10.64-27.853 10.867-31.824.641-11.237-7.678-16.888-10.664-27.992-.704-2.617-1.111-5.536-1.036-8.903.046-2.021 2.22-3.22 2.483-5.157.413-3.05-.927-6.77 0-9.563C35.699 14.207 45.439 4.316 59.034.614c3.603-.982 6.76 0 10.696 0 4.574 0 11.664-1.381 15.937 0a43.42 43.42 0 0 1 25.368 21.66c1.581 3.115-.795 6.23 0 9.64.746 3.205 4.71 6.693 4.71 10.037 0 19.65-14.277 31.511-31.925 37.688z"
        id="genericError_svg__i"
      />
      <filter
        x="-40.8%"
        y="-46.2%"
        width="181.6%"
        height="190.4%"
        filterUnits="objectBoundingBox"
        id="genericError_svg__a"
      >
        <feMorphology
          radius={1.008}
          operator="dilate"
          in="SourceAlpha"
          result="shadowSpreadOuter1"
        />
        <feOffset in="shadowSpreadOuter1" result="shadowOffsetOuter1" />
        <feGaussianBlur
          stdDeviation={7.5}
          in="shadowOffsetOuter1"
          result="shadowBlurOuter1"
        />
        <feComposite
          in="shadowBlurOuter1"
          in2="SourceAlpha"
          operator="out"
          result="shadowBlurOuter1"
        />
        <feColorMatrix
          values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0.130709135 0"
          in="shadowBlurOuter1"
        />
      </filter>
    </defs>
    <g transform="translate(2 1)" fill="none" fillRule="evenodd">
      <g transform="rotate(139 197.12 185.64)">
        <use
          fill="#000"
          filter="url(#genericError_svg__a)"
          xlinkHref="#genericError_svg__b"
        />
        <use
          stroke="#26101A"
          strokeWidth={2.016}
          fill="#F2E9E9"
          xlinkHref="#genericError_svg__b"
        />
      </g>
      <path
        stroke="#26101A"
        strokeWidth={2.24}
        fill="#5DA7B5"
        d="m317.396 141.882 39.122 9.754-9.754 39.122-39.122-9.754z"
      />
      <path
        d="M675.545 75.727c-3.193 0 3.748 3.874 2.037 10.321-1.41 5.31 1.195 7.554 3.005 8.722 6.742 4.363 11.501-1.587 11.501-1.587"
        stroke="#26101A"
        strokeWidth={1.985}
        strokeLinecap="round"
      />
      <path
        d="M676.516 75.727c-3.117 0-2.274 14.948 4.071 19.043 2.082 1.348 3.966 1.797 5.596 1.586 8.285-1.19 10.664 2.777 10.664 2.777"
        stroke="#26101A"
        strokeWidth={1.985}
        strokeLinecap="round"
      />
      <path
        d="M675.545 75.727c-.551 4.382-1.303 14.948 5.042 19.043 2.082 1.348 3.966 1.78 5.596 1.586 5.905-.793 6.698 3.968 6.698 3.968"
        stroke="#26101A"
        strokeWidth={1.985}
        strokeLinecap="round"
      />
      <path
        d="M636.38 112.062s-7.396 6.836-7.733 7.141c-3.966 3.57-7.733 9.521-7.535 39.275.043 6.348-22.407 2.777-22.407 2.777l8.328-48.796 29.347-.397z"
        stroke="#26101A"
        strokeWidth={1.588}
        fill="#4C3D6A"
        fillRule="nonzero"
      />
      <g transform="translate(543.353 20.93)">
        <mask id="genericError_svg__d" fill="#fff">
          <use xlinkHref="#genericError_svg__c" />
        </mask>
        <use
          stroke="#26101A"
          strokeWidth={1.985}
          fill="#4F728F"
          fillRule="nonzero"
          xlinkHref="#genericError_svg__c"
        />
        <path
          d="M93.027 22.446c21.305 9.066 32.66 16.725 34.063 22.976 1.404 6.25 1.404 20.58 0 42.99 10.24-31.737 15.361-47.899 15.361-48.486 0-.588-6.56-9.394-19.68-26.418L99.96 8.616l-6.932 13.83z"
          fill="#4F6486"
          mask="url(#genericError_svg__d)"
        />
        <path
          d="M135.748 56.31c1.094-13.244-2.353-23.66-10.343-31.25-7.99-7.588-17.794-11.228-29.414-10.92L88.642-2.958c16.41 4.543 28.664 9.38 36.763 14.51 8.1 5.132 18.41 14.59 30.928 28.373L135.748 56.31z"
          fill="#5DA7B5"
          mask="url(#genericError_svg__d)"
        />
      </g>
      <path
        d="M553.68 42.403S540.5 29.216 608.988 21.31c16.365-1.89 25.234 3.458 46.944 12.12 25.316 10.103 27.322 32.335 20.341 58.112-3.208 11.845-8.314 24.44-13.927 36.95-11.501 25.629-4.968 28.554-6.252 32.501-1.283 3.948-8.707 3.948-15.164 3.948-6.456 0-6.76-2.388-6.364-3.948.246-.975-1.33-13.269.812-29.328.793-5.95.793-9.918 2.776-16.662 0 0-8.443 1.829-13.484 1.587-14.277-.69-8.17 49.34-8.546 51.573-.793 4.681-7.515 5.59-13.068 4.966-5.552-.622-11.411-.039-10.221-3.157 1.88-4.931 4.868-54.176-5.047-54.176-18.243 0-.504 52.137-1.376 54.176-1.19 2.777-8.104 2.952-15.183 2.952-4.787 0-7.212-.397-6.816-4.761.179-1.984.77-8.719-2.752-16.856-5.646-13.046-5.013-28.873-5.6-36.304-.98-12.5-26.48-44.305-2.38-72.599z"
        stroke="#26101A"
        strokeWidth={1.985}
      />
      <path
        d="M638.597 139.682c1.932.36 3.507.495 5.506.135m-5.948 3.038c2.365.36 4.294.496 6.741.135"
        stroke="#26101A"
        strokeWidth={1.985}
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <g transform="rotate(-6 404.833 -4736.362)">
        <path
          d="M49.015 70.783s-2.433 6.916-7.485 12.812C32.836 93.741 15.155 94.97 16.476 97.342c.658 1.19 19.58-.347 28.15-7.546 5.377-4.516 10.842-10.71 10.842-10.71l-6.453-8.303z"
          stroke="#26101A"
          strokeWidth={1.985}
          fill="#E1D9BE"
          fillRule="nonzero"
        />
        <path
          d="M48.961 135.475c11.916-6.518 3.355-9.987 13.19-23.027.397-.523 2.035-2.44 2.8-3.268 4.55-4.94 20.331-14.388 27.192-16.364 4.988-1.43 9.955-7.83 14.451-9.918 5.485-2.546 8.715.283 12.51-2.038 1.034-.635 6.718.31 7.254-.071 5.155-3.662 6.121-9.8 8.22-12.57 25.976-34.315-24.334-60.02-49.715-50.5-25.382 9.522-37.445 90.432-37.445 90.432.15 6.759.224 11.336.218 13.732-.007 3.523.434 8.054 1.325 13.592z"
          fill="#4F597E"
          fillRule="nonzero"
          style={{
            mixBlendMode: 'darken',
          }}
        />
        <g transform="translate(87.373 .146)">
          <mask id="genericError_svg__f" fill="#fff">
            <use xlinkHref="#genericError_svg__e" />
          </mask>
          <use
            stroke="#26101A"
            strokeWidth={1.985}
            fill="#4F728F"
            fillRule="nonzero"
            xlinkHref="#genericError_svg__e"
          />
          <path
            d="M-.92 13.546c4.59-3.372 7.97-5.578 10.134-6.618 3.246-1.56 10.106-1.361 12.611-2.032 2.505-.672 12.473-1.15 16.764 1.357 4.29 2.508 14.15 7.218 18.837 24.72 3.125 11.667 4.276 2.253 3.454-28.241L17.75-10.315-3.569 3.657l2.649 9.89z"
            fill="#5DA7B5"
            mask="url(#genericError_svg__f)"
          />
          <path
            d="M.233 11.923c-.67-2.18 2.637-5.33 4.313-6.876C7.318 2.492 15.874 2.03 21.823.84c17.35-3.471 30.886 8.239 35.072 23.864C60.3 37.418 52.657 38.986 44.687 51.65s-.477 20.622-7.089 24.448c-5.022 2.906-9.033-2.074-15.775-3.661-5.794-1.365-4.648-11.303-7.932-9.918-4.445 1.868-4.414-19.836-3.569-24.597l-10.089-26z"
            stroke="#26101A"
            strokeWidth={1.985}
            mask="url(#genericError_svg__f)"
          />
        </g>
        <g transform="translate(8.334 5.235)">
          <mask id="genericError_svg__h" fill="#fff">
            <use xlinkHref="#genericError_svg__g" />
          </mask>
          <use
            stroke="#26101A"
            strokeWidth={1.985}
            fill="#4F728F"
            fillRule="nonzero"
            xlinkHref="#genericError_svg__g"
          />
          <path
            d="M46.69 2.466c-4.585-.106-8.373.242-11.364 1.043-4.487 1.203-6.476 1.911-18.47 2.512C4.863 6.622-.063 14.998-.28 21.933c-.217 6.935 1.281 6.47-3.228-4.526-3.006-7.33 1.93-12.721 14.805-16.173l35.652-6.788-.258 8.02z"
            fill="#5DA7B5"
            mask="url(#genericError_svg__h)"
          />
        </g>
        <path
          d="M66.908 19.325c.9-2.021.65-4.37-.658-6.153-2.157-2.865-7.392-5.673-12.786-7.038-9.92-2.506-25.258 1.038-32.977 3.595C11.16 12.817 8.385 15.76 8.433 24.462c.052 9.364 8.586 16.238 8.784 24.272.189 7.611.733 11.259 6.293 23.49 5.559 12.232 4.641-4.74 18.818-5.044 5.45-.116 2.503-1.936 5.29-.166 3.768 2.38 6.715-18.832 6.6-23.577l12.69-24.112z"
          stroke="#26101A"
          strokeWidth={1.985}
        />
        <path
          d="M84.013 87.064c-6.345 2.217-19.88 15.57-24.786 18.844-9.518 6.347-7.236 29.731-25.83 27.513-19.216-2.293-23.962-12.394-26.685-19.677-.922-2.467-8.657-18.437-3.553-16.15 5.103 2.287 4.782-5.55 9.037-1.003 1.968 2.101-.385 4.36.702 7.57 2.442 7.22 7.893 15.99 16.861 15.452 6.864-.412 10.8-29.371 11.027-33.343.641-11.237-7.678-16.887-10.664-27.991-.704-2.617-1.111-5.536-1.036-8.903.046-2.022 2.22-3.22 2.483-5.157.413-3.05-.927-6.77 0-9.563C35.892 21.632 45.632 11.741 59.227 8.04c3.603-.982 6.76 0 10.695 0 4.575 0 11.665-1.381 15.938 0a43.42 43.42 0 0 1 25.368 21.66c1.581 3.114-.795 6.23 0 9.64.746 3.205 4.71 6.693 4.71 10.037 0 19.65-14.277 31.511-31.925 37.688z"
          fill="#4F728F"
          fillRule="nonzero"
        />
        <path
          d="M61.606 7.522S46.14 10.82 46.14 35.292c0 14.282 15.046 17.34 11.105 30.548-.607 2.033-12.663 6.668-12.108 8.617.384 1.347 9.4-1.657 12.738.502 3.339 2.16 4.306 3.647 6.298 5.757.372.394-7.423-4.17-12.492-4.165-5.493.005-8.222 4.562-5.75 4.562 1.67 0 3.498-1.65 5.75-1.59 2.77.076 6.008 1.928 9.716 5.557-5.41-1.096-8.793-1.625-10.149-1.587-5.375.15-5.291 2.184-5.714 2.777-.822 1.152 7.919.098 7.931 1.587.031 3.644-6.373 11.766-8.614 14.2-3.738 4.06-6.617 22.327-16.281 20.154-8.542-1.92-13.264-7.407-14.168-16.46 7.622 10.443 13.402 14.24 17.339 11.39 6.578-4.762 10.715-33.107 6.071-43.764-1.451-3.336-5.5-7.656-6.94-13.29 0 0-3.432-7.005-1.785-10.513 3.353-7.136 2.199-15.708 7.007-23.339C47.593 7.987 61.606 7.522 61.606 7.522z"
          fill="#4F6486"
          fillRule="nonzero"
        />
        <g transform="translate(.193 7.425)">
          <mask id="genericError_svg__j" fill="#fff">
            <use xlinkHref="#genericError_svg__i" />
          </mask>
          <path
            d="M5.011 81.27c-2.605 4.429-3.17 7.434-1.697 9.016 2.21 2.373 6.144-4.162 6.51 3.494.367 7.655 15.479 28.016 23.672 18.82 5.462-6.13 8.031-18.81 7.707-38.04L28.51 92.935l-8.457-4.748L5.01 81.27zM44.116 7.618c11.27-5.184 17.278-6.384 18.024-3.6 1.12 4.178-3.425 2.334-1.771 3.1 1.654.765 12.914-7.68 25.099-3.808 8.122 2.581 14.99 6.553 20.603 11.915l13.915-10.909L75.442-18.15l-13.67 7.754L44.116 7.618z"
            fill="#5DA7B5"
            mask="url(#genericError_svg__j)"
          />
        </g>
        <path
          d="M49.709 62.07c0-1.642-1.333-2.578-2.975-2.578-1.642 0-2.974.936-2.974 2.579a2.975 2.975 0 1 0 5.949 0zm34.776.112c.635-2.71-1.586-4.245-4.109-3.805a2.979 2.979 0 0 0-2.502 2.543c-.262 2.444 2.24 5.181 4.838 3.714a3.785 3.785 0 0 0 1.773-2.452z"
          fill="#26101A"
          fillRule="nonzero"
        />
        <path
          d="M45.136 87.46s2.727-1.519 6.03-1.365c3.302.154 5.92 1.838 7.057 3.746M7.98 101.834c-2.379 1.487-4.222.654-4.222.654M45.533 81.51c2.03-1.344 4.451-1.856 7.265-1.537 4.785.543 8.085 3.035 8.878 5.018M45.93 75.56c1.79-1.44 4.542-1.905 8.257-1.396 5.574.764 9.334 4.52 9.985 6.95M84.043 50.72c-.57 1.106-.47 2.138-.001 3.065.95 1.874 3.413 3.318 4.923 4.073 1.224.612 2.519.842 3.326.625m-46.95-6.885c1.946-1.315 3.066.06 4.711.756 1.925.816 4.017-.866 4.71 1.725m-43.607 62.044c-.816.686-2.442 1.27-3.877 1.096m7.585 2.017c-.816.686-2.442 1.27-3.877 1.096m8.351 1.812c-.816.686-2.442 1.27-3.877 1.096M33.7 56.65c-2.214 3.473-.661 8.396.612 11.768M15.55 20.436c0-4.562 5.646-8.55 12.862-8.55 11.831 0 17.216-6.76 28.225-3.372"
          stroke="#26101A"
          strokeWidth={1.985}
          strokeLinecap="round"
        />
        <path
          d="M100.635 77.528s6.037-4.728 9.64-4.728c2.816 0 6.659-14.131 16.166-23.246 9.683-9.282 11.745-25.6 3.978-21.943-4.163 1.96-7.587 3.745-10.442 3.621-2.642-.114-3.898-3.549-8.986-10.276-5.089-6.727-8.222-1.076-.54 8.942 1.32 1.722 1.98 14.141 2.217 16.661.781 8.306-1.566 6.596-2.392 10.979-2.512 13.326-9.641 19.99-9.641 19.99z"
          fill="#4E567C"
          fillRule="nonzero"
          style={{
            mixBlendMode: 'darken',
          }}
        />
        <path
          d="M7.489 102.347c3.09 1.318-5.813 1.518 1.2 8.95 7.013 7.433 12.031 17.795 27.863 16.776 6.066-.392 12.135-16.698 14.032-19.068 8.07-10.08 9.6-9.918 17.367-14.999a53.134 53.134 0 0 1 14.674-6.744s-7.62 5.498-11.637 10.258c-2.205 2.613-6.292 5.26-15.184 13.568-7.825 7.311-4.994 17.265-17.244 21.217-4.663 1.505-11.796-.584-20.612-4.4-8.816-3.817-11.605-11.618-14.357-17.143l-2.77-10.82c.184-1.755 1.332-2.521 3.442-2.3 3.491.368-1.121 2.852 3.226 4.705z"
          fill="#4D4C75"
          fillRule="nonzero"
          style={{
            mixBlendMode: 'darken',
          }}
        />
        <ellipse
          fill="#C3E5D7"
          fillRule="nonzero"
          cx={79.471}
          cy={62.405}
          rx={1}
          ry={1.19}
        />
        <circle
          fill="#C3E5D7"
          fillRule="nonzero"
          cx={45.166}
          cy={63.079}
          r={1}
        />
        <path
          d="M35.035 21.804c-4.118 5.089-6.756 9.141-8.457 14.237-.649 1.943-.23 4.44-.636 6.801-.273 1.585-1.717 2.85-1.615 4.749.968 17.999 5.377 27.667 5.377 27.667l4.769-5.105-.794-2.965s-5.386-9.68-4.593-17.614c.227-2.266 1.255-2.737 1.74-5.355 1.215-6.547 2.736-15.727 6.985-20.828A86.385 86.385 0 0 1 50.105 11.49c-.944.138-9.593 3.542-15.07 10.314z"
          fill="#4F597E"
          fillRule="nonzero"
          style={{
            mixBlendMode: 'darken',
          }}
        />
        <path
          d="M50.105 11.49c-.944.138-10.31 1.586-16.656 5.157-9.34 5.252-15.98 12.222-13.88 23.406 2.605 13.882 6.503 36.178 6.503 36.178l3.632-.973L27.5 67.03s-3.173-10.712-2.38-18.646c.794-7.934 3.966-17.455 9.915-24.596 5.73-6.895 15.07-12.298 15.07-12.298z"
          fill="#4F6486"
          fillRule="nonzero"
          style={{
            mixBlendMode: 'darken',
          }}
        />
        <path
          d="M84.013 87.064c-6.345 2.217-19.865 17.811-24.77 21.084-9.519 6.348-9.926 26.946-25.635 24.936-19.198-2.457-26.202-13.353-28.925-20.637-.922-2.466-6.628-17.14-1.524-14.853 2.497 1.12 4.917-2.103 5.553-2.734.637-.63 2.5-1.736 3.842 1.23 2.334 5.158-1.647 10.397 11.467 20.239.523.393 1.03.713 1.524.964 11.857 6.061 15.024-27.21 15.24-31.023.642-11.237-7.677-16.887-10.663-27.991-.704-2.617-1.111-5.536-1.036-8.903.046-2.022 2.22-3.22 2.483-5.157.413-3.05-.927-6.77 0-9.563C35.892 21.632 45.632 11.741 59.227 8.04c3.603-.982 6.76 0 10.695 0 4.575 0 11.665-1.381 15.938 0a43.42 43.42 0 0 1 25.368 21.66c1.581 3.114-.795 6.23 0 9.64.746 3.205 4.71 6.693 4.71 10.037 0 19.65-14.277 31.511-31.925 37.688z"
          stroke="#26101A"
          strokeWidth={1.985}
        />
        <path
          d="M82.443 83.001c2.137.397 2.38 2.543 2.01 4.678-.587 3.265-2.82 10.172-5.496 14.654-5.295 8.86-23.6 22.004-25.56 20.207-1.19-1.075 15.784-16.577 18.366-21.56 4.04-7.803 4.584-16.019 4.584-16.019s.944-2.634 5.29-2.083c.286.056.556.08.806.123z"
          stroke="#26101A"
          strokeWidth={1.985}
          fill="#E1D9BE"
          fillRule="nonzero"
          strokeLinecap="round"
          strokeLinejoin="round"
        />
        <path
          stroke="#26101A"
          strokeWidth={1.985}
          strokeLinecap="round"
          d="m75.673 79.526.397 6.744M59.623 15.89s3.57-7.935 15.467-7.935m12.48 72.761-1.586 5.158"
        />
      </g>
      <path
        d="M560.47 148.41a7.061 7.061 0 0 0 5.505-1.274m-4.315 4.05a7.061 7.061 0 0 0 5.505-1.273m32.916-4.761a7.481 7.481 0 0 0 5.949.397m11.5-21.422c1.587-8.728 4.922-16.266 5.318-23.803m12.529 30.15c1.19-20.63 10.708-21.82 7.138-36.895m-42.434 55.144s3.605 1.741 6.742.397"
        stroke="#26101A"
        strokeWidth={1.985}
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M616.857.671c1.844-.164 3.616-.027 5.315.411 2.549.658 6.209 2.3 6.815 2.876"
        stroke="#26101A"
        strokeLinecap="round"
        strokeWidth={1.12}
      />
      <g stroke="#26101A" strokeLinejoin="round" strokeWidth={2.24}>
        <path
          fill="#A5597E"
          d="m299.893 173.802 25.616-6.027 8.386 25.666-26.388 5.229z"
        />
        <path
          fill="#D06868"
          d="m325.51 167.775 8.512-14.618 7.228 24.468-7.355 15.816z"
        />
        <path
          fill="#F59178"
          d="m309.19 159.199 24.832-6.042-8.513 14.618-25.616 6.027z"
        />
      </g>
      <g transform="rotate(3 -2224.575 7875.563)" stroke="#26101A">
        <path
          strokeWidth={2.24}
          fill="#A5597E"
          strokeLinejoin="round"
          d="m41.478 1.553 82.536 1.033L85.23 23.393.254 21.983z"
        />
        <path
          strokeWidth={2.24}
          fill="#623578"
          strokeLinejoin="round"
          d="M.373 22.649 41.478 1.553l3.885 49.354L3.338 68.77z"
        />
        <circle
          strokeWidth={2.24}
          fill="#65ADBB"
          cx={70.398}
          cy={20.156}
          r={13.16}
        />
        <path
          strokeWidth={2.24}
          fill="#65ADBB"
          d="m54.489 7.336.189 30.198L28.78 23.18z"
        />
        <circle strokeWidth={2.24} cx={94.545} cy={30.002} r={22.12} />
        <circle
          strokeWidth={2.52}
          fill="#623578"
          cx={80.816}
          cy={26.339}
          r={13.16}
        />
        <g strokeLinejoin="round">
          <path
            strokeWidth={2.24}
            fill="#A5597E"
            d="m51.382 16.34 16.153 9.614-8.426 17.349-16.153-10.407z"
          />
          <path
            strokeWidth={1.985}
            fill="#D06868"
            d="m67.535 25.954 11.627-3.287-8.426 16.159-11.627 4.477z"
          />
          <path
            strokeWidth={2.24}
            fill="#F59178"
            d="m63.406 13.45 15.756 9.217-11.627 3.287-16.153-9.614z"
          />
        </g>
        <path
          strokeWidth={2.24}
          fill="#D06868"
          strokeLinejoin="round"
          d="M85.617 23.931 124.4 2.65l2.097 45.554L90.323 72.09z"
        />
        <path
          strokeWidth={2.24}
          fill="#F59178"
          strokeLinejoin="round"
          d="m.373 22.649 84.898 1.433 4.702 47.847L3.338 68.77z"
        />
      </g>
      <g stroke="#26101A" strokeLinejoin="round" strokeWidth={1.985}>
        <path
          fill="#A5597E"
          d="m184.767 165.037 16.153 9.613-8.427 17.35-16.153-10.407z"
        />
        <path
          fill="#D06868"
          d="m200.92 174.65 11.627-3.287-8.427 16.159-11.627 4.477z"
        />
        <path
          fill="#F59178"
          d="m196.79 162.147 15.757 9.216-11.627 3.287-16.153-9.613z"
        />
      </g>
      <path
        stroke="#26101A"
        strokeWidth={2.52}
        fill="#623578"
        d="m272.812 141.678 12.79 5.694-5.694 12.79-12.79-5.694z"
      />
      <circle
        stroke="#26101A"
        strokeWidth={2.24}
        fill="#D16F6F"
        cx={175.56}
        cy={189.56}
        r={13.16}
      />
      <circle
        stroke="#26101A"
        strokeWidth={2.268}
        cx={237.72}
        cy={162.12}
        r={22.12}
      />
      <circle
        stroke="#26101A"
        strokeWidth={2.24}
        fill="#623578"
        cx={196.28}
        cy={195.16}
        r={13.16}
      />
      <g stroke="#26101A" strokeLinejoin="round" strokeWidth={1.985}>
        <path
          fill="#A5597E"
          d="m297.887 132.557 16.153 9.613-8.427 17.349-16.153-10.407z"
        />
        <path
          fill="#D06868"
          d="m314.04 142.17 11.627-3.287-8.427 16.159-11.627 4.477z"
        />
        <path
          fill="#F59178"
          d="m309.91 129.667 15.757 9.216-11.627 3.287-16.153-9.613z"
        />
      </g>
      <path
        stroke="#26101A"
        strokeWidth={2.24}
        fill="#5DA7B5"
        d="m214.795 145.758 30.887 25.917-25.917 30.887-30.887-25.917z"
      />
      <circle
        stroke="#26101A"
        strokeWidth={2.261}
        cx={262.92}
        cy={154.84}
        r={38.36}
      />
      <g stroke="#26101A">
        <path
          d="M487.927 79.605c6.644-5.175 17.066-2.944 23.28 4.984 6.212 7.927 5.862 18.55-.782 23.724-6.644 5.175-17.067 2.944-23.28-4.984-6.212-7.927-5.863-18.55.782-23.724zm1.616 3.422c-4.728 3.682-4.812 12.229-.015 18.35 4.798 6.122 12.66 7.5 17.387 3.817 4.727-3.682 4.67-11.63-.126-17.75-4.798-6.122-12.519-8.1-17.246-4.417z"
          strokeWidth={2.24}
          fill="#F59178"
        />
        <path
          d="M509.938 106.294a2.82 2.82 0 0 1 3.893.493l2.184 2.727c-.85 1.03-1.039 1.92-.567 2.668.472.749 1.675 1.468 3.609 2.156l3.52 4.565a2.09 2.09 0 0 1-.371 2.924 2.145 2.145 0 0 1-2.968-.322l-6.41-7.726-2.203-2.654-1.046-1.163a2.498 2.498 0 0 1 .36-3.668z"
          strokeWidth={2.24}
          fill="#A5597E"
          strokeLinejoin="round"
        />
        <path
          d="M489.387 92.386c.94 1.825 2.25 3.304 3.861 4.632 2.924 2.41 6.146 3.757 10.731 3.393a18.76 18.76 0 0 0 1.343-.155c-.124.615-.344 1.122-.587 1.472-.459.66-1.474 1.439-2.903 1.766-1.005.23-2.216.242-3.553-.18-1.889-.598-4.07-2.078-6.251-5.162-1.554-2.198-2.341-4.124-2.64-5.766z"
          strokeWidth={1.68}
          fill="#C3E5D7"
          strokeLinejoin="round"
        />
      </g>
      <path
        stroke="#26101A"
        strokeWidth={2.24}
        fill="#D16F6F"
        d="m303.282 141.603 17.675 36.24-36.24 17.674-17.674-36.24z"
      />
      <circle
        stroke="#26101A"
        strokeWidth={2.24}
        fill="#65ADBB"
        cx={265.64}
        cy={163.192}
        r={13.16}
      />
      <circle
        stroke="#26101A"
        strokeWidth={2.24}
        fill="#65ADBB"
        cx={429.8}
        cy={196.84}
        r={13.16}
      />
      <path
        stroke="#26101A"
        strokeWidth={2.24}
        fill="#D16F6F"
        d="m204.006 174.607 39.707-7.001 7.001 39.707-39.707 7.001z"
      />
      <circle
        stroke="#26101A"
        strokeWidth={2.24}
        fill="#623578"
        cx={280}
        cy={184.8}
        r={16.24}
      />
      <path
        stroke="#26101A"
        strokeWidth={2.52}
        fill="#65ADBB"
        d="m287.167 108.523-1.654 30.309-25.361-16.47zm47.6 71.12-1.654 30.309-25.361-16.47z"
      />
      <circle
        stroke="#26101A"
        strokeWidth={2.24}
        fill="#65ADBB"
        cx={296.52}
        cy={192.36}
        r={13.16}
      />
      <path
        stroke="#26101A"
        strokeWidth={2.52}
        fill="#623578"
        d="m245.36 159.6 20.56.23-10.08 17.46z"
      />
      <g stroke="#26101A" strokeLinejoin="round" strokeWidth={1.985}>
        <path
          fill="#A5597E"
          d="m232.927 166.157 16.153 9.613-8.427 17.349-16.153-10.407z"
        />
        <path
          fill="#D06868"
          d="m249.08 175.77 11.627-3.287-8.427 16.159-11.627 4.477z"
        />
        <path
          fill="#F59178"
          d="m244.95 163.267 15.757 9.216-11.627 3.287-16.153-9.613z"
        />
      </g>
      <circle
        stroke="#26101A"
        strokeWidth={2.24}
        fill="#623578"
        cx={311.08}
        cy={199.64}
        r={13.16}
      />
      <circle
        stroke="#26101A"
        strokeWidth={2.24}
        fill="#65ADBB"
        cx={250.88}
        cy={198.24}
        r={16.8}
      />
      <g stroke="#26101A" strokeLinejoin="round" strokeWidth={2.24}>
        <path
          fill="#A5597E"
          d="m261.962 196.564 26.226 2.183.044 27.002-26.713-3.182z"
        />
        <path
          fill="#D06868"
          d="m288.188 198.747 12.613-11.27-.686 25.503-11.883 12.769z"
        />
        <path
          fill="#F59178"
          d="m275.318 185.548 25.483 1.928-12.613 11.271-26.226-2.183z"
        />
      </g>
      <path
        stroke="#26101A"
        strokeWidth={2.24}
        fill="#623578"
        d="M142.52 208.32 154 228.48h-22.96z"
      />
      <path
        d="M607.04 74.706c1.237 0 2.24-1.12 2.24-2.501 0-.921-.747-2.7-2.24-5.34-1.493 2.64-2.24 4.419-2.24 5.34 0 1.381 1.003 2.5 2.24 2.5zm-5.6-11.2c1.237 0 2.24-1.12 2.24-2.501 0-.921-.747-2.7-2.24-5.34-1.493 2.64-2.24 4.419-2.24 5.34 0 1.381 1.003 2.5 2.24 2.5z"
        stroke="#26101A"
        strokeWidth={1.4}
        fill="#C3E5D7"
      />
    </g>
  </svg>
);

export default SvgGenericError;
