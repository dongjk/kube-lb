# Nginx proxy for kubernetes loadbalancing
Inspired by [kubernetes ingress controller](https://github.com/kubernetes/contrib/tree/master/ingress/controllers/nginx-alpha)

This is an external loadbalancer for kubernetes, using Nginx.

![Image of external loadbalancer kubernetes]  (img/Picture2.png)
## How to use

1. Having a external node, install flannel and set it as a subnet of kubernetes cluster flannel.
2. get kubernetes ca.cert and user token.
3. run the image in this node
    ```
    export CA_PATH=<path of ca.crt, token>
    export MASTER_IP=<k8s apiserver IP>
    export MATER_PORT=<k8s apiserver port, like 6443>

    docker run -d -p 80:80 -v "$CA_PATH":"/var/run/secrets/kubernetes.io/serviceaccount/" -e "KUBERNETES_SERVICE_HOST=$MASTER_IP" -e "KUBERNETES_SERVICE_PORT=$MASTER_PORT" lordx/kube-lb:0.1
    ```
4. Access endpoints in this url ```http://<nginx node ip>/<RC name>/```


