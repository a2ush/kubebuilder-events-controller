# kubebuilder-events-controller

This controller puts kubernetes events to CloudWatch Logs.

On Terminal
```
$ kubectl get event -A
NAMESPACE     LAST SEEN   TYPE      REASON             OBJECT                         MESSAGE
kube-system   100s        Warning   FailedScheduling   pod/coredns-659f9d44fd-gnk2v   no nodes available to schedule pods
kube-system   2m10s       Warning   FailedScheduling   pod/coredns-659f9d44fd-lsjhs   no nodes available to schedule pods
...
```

In CloudWatch Logs
![image](images/console.png)

## Setting

You can set environment values. [[manifest]](config/manager/manager.yaml)
```
        env:
        - name: CW_LOG_GROUP_NAME
          value: "/kubernetes/event-log-group"
          name: CW_LOG_STREAM_NAME
          value: "kubernetes-event-log-stream"          
          name: AWS_REGION
          value: "ap-northeast-1"      
```

* CW_LOG_GROUP_NAME: CloudWatch Logs group name (default - `/kubernetes/event-log-group`)
* CW_LOG_STREAM_NAME: CloudWatch Logs stream name (default - `kubernetes-event-log-stream`)
* AWS_REGION: region (default - `ap-northeast-1`)

And when you want to use IRSA, you should perform the following command.
```
$ ACCOUNT_ID=$(aws sts get-caller-identity --query "Account" --output text)
$ OIDC_PROVIDER=$(aws eks describe-cluster --name cluster-name --query "cluster.identity.oidc.issuer" --output text | sed -e "s/^https:\/\///")
$ aws iam create-role --role-name IAM_ROLE_NAME --assume-role-policy-document file://trust.json --description "IAM_ROLE_DESCRIPTION"
``` 


## How to deploy this controller as a pod in your cluster
```
$ git clone https://github.com/a2ush/kubebuilder-events-controller.git
$ cd kubebuilder-events-controller
$ make docker-build docker-push IMG=<registry>/<project-name>:tag
$ make deploy IMG=<registry>/<project-name>:tag
```

If you want to use IRSA, you should perform the following command.
```
$ git clone https://github.com/a2ush/kubebuilder-events-controller.git
$ cd kubebuilder-events-controller
$ make docker-build docker-push IMG=<registry>/<project-name>:tag

$ ACCOUNT_ID=$(aws sts get-caller-identity --query "Account" --output text)
$ OIDC_PROVIDER=$(aws eks describe-cluster --name cluster-name --query "cluster.identity.oidc.issuer" --output text | sed -e "s/^https:\/\///")
$ sed -i -e "s/<ACCOUNT_ID>/$ACCOUNT_ID/" -e "s|<OIDC_PROVIDER>|${OIDC_PROVIDER}|" irsa/trust.json 
$ aws iam create-role --role-name kubebuilder-events-controller --assume-role-policy-document file://irsa/trust.json --description "For kubebuilder-events-controller role"
$ aws iam attach-role-policy --role-name kubebuilder-events-controller --policy-arn=arn:aws:iam::aws:policy/CloudWatchAgentServerPolicy

$ make deploy IMG=<registry>/<project-name>:tag AWS_ACCOUNT_ID=${ACCOUNT_ID}
```

Ex)
```
$ kubectl version --short
Client Version: v1.20.4-eks-6b7464
Server Version: v1.21.5-eks-bc4871b

$ aws ecr create-repository --repository-name kubebuilder-events-controller
$ aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin 111122223333.dkr.ecr.ap-northeast-1.amazonaws.com

$ make docker-build docker-push IMG=111122223333.dkr.ecr.ap-northeast-1.amazonaws.com/kubebuilder-events-controller:latest
$ make deploy IMG=111122223333.dkr.ecr.ap-northeast-1.amazonaws.com/kubebuilder-events-controller:latest
```

Environment
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

## How to test without deploying
```
$ git clone https://github.com/a2ush/kubebuilder-events-controller.git
$ cd kubebuilder-events-controller
$ make
$ make run
```
