# kubebuilder-events-controller

This controller puts kubernetes events to CloudWatch Logs.


## Setting

You can set environment values.
* CW_LOG_GROUP_NAME: CloudWatch Logs group name (default - /kubernetes/event-log-group)
* CW_LOG_STREAM_NAME: CloudWatch Logs stream name (default - kubernetes-event-log-stream)
* AWS_REGION: region (default - ap-northeast-1)

## How to deploy 
```
make docker-build docker-push IMG=<registry>/<project-name>:tag
make deploy IMG=<registry>/<project-name>:tag
```

```
$ kubectl get all -n kubebuilder-events-controller-system 
NAME                                                                  READY   STATUS    RESTARTS   AGE
pod/kubebuilder-events-controller-controller-manager-5665ff56994wwm   2/2     Running   0          17m

NAME                                                       TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
service/kubebuilder-events-controller-cm-metrics-service   ClusterIP   10.100.169.42   <none>        8443/TCP   17m

NAME                                                               READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/kubebuilder-events-controller-controller-manager   1/1     1            1           17m

NAME                                                                          DESIRED   CURRENT   READY   AGE
replicaset.apps/kubebuilder-events-controller-controller-manager-5665ff5695   1         1         1       17m
```
