# Authorization Webhook for use with role-keeper

Run apiserver with the following args: `--kube-apiserver-arg authorization-mode=Webhook,Node,RBAC --kube-apiserver-arg authorization-webhook-config-file=<path-to-webhook-config-yaml>`

An example of the webhook config yaml reference above can be found in this project at `config/`
