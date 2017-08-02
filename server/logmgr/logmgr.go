package logmgr

import (
	"bufio"
	"io"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/arlert/malcolm/model"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
)

type LogMgr struct {
	client *kubernetes.Clientset
}

func NewLogMgr(client *kubernetes.Clientset) *LogMgr {
	logMgr := &LogMgr{client: client}
	return logMgr
}

func (l *LogMgr) GetLog(build *model.Build, writer io.Writer) error {
	for step_index, step := range build.Steps {
		if step_index != 0 {
			time.Sleep(5 * 1e9)
		}
		job := step.K8sjob
		// set pod name is job is ignored
		podlist, err := l.client.Pods(job.Namespace).List(meta_v1.ListOptions{
			LabelSelector: "job-name=" + job.Name,
		})
		logrus.Debug("podlist ", podlist)
		if err != nil {
			logrus.Debug("client.Pods ", err)
			continue
		}
		if len(podlist.Items) == 0 {
			continue
		}
		pod := podlist.Items[0]
		for index, container := range pod.Spec.InitContainers {
			if index != 0 {
				time.Sleep(5 * 1e9)
			}
			req := l.client.Pods(job.Namespace).GetLogs(pod.Name,
				&v1.PodLogOptions{Container: container.Name, Follow: true})
			if req == nil {
				logrus.Debug("GetLogs req nil")
				continue
			}
			reader, err := req.Stream()
			defer reader.Close()
			if err != nil {
				logrus.Debug("req.Stream error", err)
				continue
			}
			logrus.Debug("reading initcontainer logs")
			_, err = io.Copy(writer, reader)
			if err == io.EOF {
				logrus.Debug("eof")
				continue
			}
		}

		var readers []io.Reader
		for index, container := range pod.Spec.Containers {
			if index != 0 {
				time.Sleep(5 * 1e9)
			}
			req := l.client.Pods(step.K8sjob.Namespace).GetLogs(pod.Name,
				&v1.PodLogOptions{Container: container.Name, Follow: true})
			if req == nil {
				logrus.Debug("GetLogs req nil")
				continue
			}
			reader, err := req.Stream()
			if err != nil {
				logrus.Debug("req.Stream error", err)
				continue
			} else {
				defer reader.Close()
			}
			readers = append(readers, reader)
		}
		readv(readers, writer)
	}
	return nil
}

type Msg struct {
	data []byte
	err  error
}

func readv(readers []io.Reader, writer io.Writer) {
	buf := make(chan *Msg)
	counter := len(readers)
	if counter == 0 {
		return
	}
	for _, reader := range readers {
		go func() {
			logrus.Debug("reading container logs")
			bufreader := bufio.NewReader(reader)
			for {
				data, err := bufreader.ReadBytes('\n')
				buf <- &Msg{
					data: data,
					err:  err,
				}
				if err != nil {
					return
				}
			}
		}()
	}
	for {
		select {
		case msg := <-buf:
			if msg.err != nil {
				counter -= 1
				if counter == 0 {
					return
				}
			} else {
				_, err := writer.Write(msg.data)
				if err != nil {
					logrus.Debug("writing container logs error", err)
					return
				}
			}
		case <-time.After(time.Minute * 60):
			logrus.Debug("timeout reading logs")
			return

		}
	}
}
