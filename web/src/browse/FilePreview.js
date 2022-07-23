import { QueryClient, QueryClientProvider, useQuery } from '@tanstack/react-query'
import MonacoEditor from 'react-monaco-editor';
import { useEffect, useRef } from 'react';

function FilePreview(props) {
  let { isLoading, error, data } = useQuery(['searchResults', props.project, props.file], () =>
      fetch(`/api/${props.project}/file/${props.file}`).then(res =>
          res.text()
      )
  )

  if (isLoading) return 'Loading...'
  if (error) return 'An error has occurred: ' + error.message
  return (
      <div style={{height: "100%"}}>
        <MonacoEditor
          language="java"
          theme="vs-dark"
          value={data}
          height="100%"
        />
      </div>
  )
}

export default FilePreview;