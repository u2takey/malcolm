package pipemgr

var VolumnTemplateDefault = `
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  namespace: {{ .Meta.Namespace }}
  name: {{ .Meta.BuildID }}
spec:
  storageClassName: {{ .Pipe.StorageClass }}
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: {{ .Pipe.StorageSize }}
`

var JobTemplateDefault = `
{{ $data := . }}
apiVersion: batch/v1
kind: Job
metadata:
  name: {{ .Meta.BuildID }}-{{ .Meta.TaskID }}
  namespace: {{ .Meta.Namespace }}
  labels:
    malcolm: malcolm-job
    pipe: {{ .Meta.PipeID }}
    build: {{ .Meta.BuildID }}
    task: {{ .Meta.TaskID }}
spec:
  completions: 1
  parallelism: 1
  template:
    metadata:
      name: {{ .TaskGroup.Title | str2title }}
      labels:
        type: {{ .Meta.Type }}
        pipe: {{ .Meta.PipeID }}
        build: {{ .Meta.BuildID }}
        task: {{ .Meta.TaskID }}
        malcolm: malcolm-job
    spec:
{{ range $Index, $Task := .TaskGroup.PreTasks }}
      initContainers:
      - name: {{ $Index }}-{{ $Task.Title | str2title }}
        image: {{ $Task.Plugin }}
        imagePullPolicy: {{ $Task.PullPolicy }}
        args:
  {{ range $Task.Args }}
        - {{.}}
  {{ end }}
        command:
  {{ range $Task.Command }}
        - {{.}}
  {{ end }}
        env:
  {{ range $Key, $Val := $Task.Environment }}
        - name: {{ $Key }}
          value: {{ $Val }}
  {{ end }} 
  {{ range $Key, $Val := $Task.Environment }}
        - name: {{ $Key }}
          valueFrom:
            secretKeyRef:
              name: {{ $Val }}
              key: data  
  {{ end }}          
{{ end }}
{{ range $Index, $Task := .TaskGroup.Tasks }}
      containers:
      - name: {{ $Index }}-{{ $Task.Title | str2title }}
        image: {{ $Task.Plugin }}
        imagePullPolicy: {{ $Task.PullPolicy }}
        args:
  {{ range $Task.Args }}
        - {{.}}
  {{ end }}
        command:
  {{ range $Task.Command }}
        - {{.}}
  {{ end }}
        env:
  {{ range $Key, $Val := $Task.Environment }}
        - name: {{ $Key }}
          value: {{ $Val }}
  {{ end }}
  {{ range $Key, $Val := $Task.Environment }}
        - name: {{ $Key }}
          valueFrom:
            secretKeyRef:
              name: {{ $Val }}
              key: data  
  {{ end }}
  {{ if $data.Pipe.StorageClass }}
        volumeMounts:
        - name: workspace
          mountPath: {{ $data.Pipe.WorkSpace }}
  {{ end }}
        WorkingDir: {{ $data.Pipe.WorkSpace }}
{{ end }}
      restartPolicy: Never
{{ if .Pipe.StorageClass }}
      volumes:
      - name: workspace
        persistentVolumeClaim:
          claimName: {{ .Meta.BuildID  }}
{{ end }}

`

var ServiceTemplateDefault = `
apiVersion: batch/v1
kind: Deployment
metadata:
  name: {{ .Meta.BuildID }}-{{ .Meta.TaskID }}-{{ .Meta.ServiceID }}
  namespace: {{.Meta.Namespace}}
  labels:
    pipe: {{ .Meta.PipeID }}
    build: {{ .Meta.BuildID }}
    task: {{ .Meta.TaskID }}
    malcolm: malcolm-service
spec:
  replicas: 1
  template:
    metadata:
      name: {{ .Task.Title | str2title }}
      id: {{ .Meta.BuildID }}-{{ .Meta.TaskID }}-{{ .Meta.ServiceID }}
      labels:
        pipe: {{ .Meta.PipeID }}
        build: {{ .Meta.BuildID }}
        task: {{ .Meta.TaskID }}
        type: {{ .Meta.Type }}
    spec:
        containers:
      - name: {{ $Index }}-{{ Task.Title | str2title }}
        image: {{ $Task.Plugin }}
        imagePullPolicy: {{ .Task.PullPolicy }}
        ports:
  {{ range .Task.Ports }}
        - containerPort: {{.}}
  {{ end }}
        args:
  {{ range .Task.Args }}
        - {{.}}
  {{ end }}
        command:
  {{ range .Task.Command }}
        - {{.}}
  {{ end }}
        env:
  {{ range $Key, $Val := .Task.Environment }}
        - name: {{ $Key }}
          value: {{ $Val }}
  {{ end }}
  {{ range $Key, $Val := .Task.Environment }}
        - name: {{ $Key }}
          valueFrom:
            secretKeyRef:
              name: {{ $Val }}
              key: data  
  {{ end }}
{{ end }}
---
apiVersion: v1
kind: Service
metadata:
  name: {{ .Meta.BuildID }}-{{ .Meta.TaskID }}-{{ .Meta.ServiceID }}
  namespace: {{.Meta.Namespace}}
spec:
  ports:
{{ range $Index, $Val := .Task.Ports }}
    - name: {{ .Index }}
      port: {{ .Val }}
{{ end }}
  selector:
    id: {{ .Meta.BuildID }}-{{ .Meta.TaskID }}-{{ .Meta.ServiceID }}
`
