import { QueryClient, QueryClientProvider, useQuery } from '@tanstack/react-query'
import MonacoEditor from 'react-monaco-editor';

function GuessLanguage(name) {
  const fileExt = name.toLowerCase().split('.').pop();
  switch (fileExt) {
    case "css":
      return "css"
    case "java":
      return "java"
    case "c":
    case "h":
      return "c"
    case "cpp":
    case "hpp":
    case "cxx":
    case "hxx":
      return "cpp"
    case "md":
      return "markdown"
    case "cmd":
      return "bat"
    case "py":
      return "python"
    default:
      return fileExt
  }
}

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
          language={GuessLanguage(props.file)}
          theme="vs-dark"
          value={data}
          height="100%"
        />
      </div>
  )
}

export default FilePreview;