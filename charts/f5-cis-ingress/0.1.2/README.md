# Helm Chart for Managing Ingress Resources with a BIG-IP Device

This chart simplifies repeatable, versioned use of the [F5 BIG-IP Controller as an Ingress Controller](http://clouddocs.f5.com/containers/latest/kubernetes/kctlr-k8s-ingress-ctlr.html) in Kubernetes. 

### Prereqisites

- Install [Helm with Tiller](https://docs.helm.sh/using_helm/#installing-helm) on your cluster with appropriate permissions.
- Deploy the F5 BIG-IP Controller in your cluster. You can use the [f5-bigip-ctlr chart](https://github.com/F5Networks/charts/tree/master/src/stable/f5-bigip-ctlr) to deploy the Controller or you can deploy it [manually](http://clouddocs.f5.com/containers/latest/kubernetes/kctlr-app-install.html). 
- Deploy the Pods/Services accepting traffic from the Ingress.

> **Note:** This chart and the [f5-bigip-controller](https://github.com/recursivelycurious/charts/tree/wip/src/stable/f5-bigip-ctlr) chart can be used *independently or together*.  
> If you or your organization author your own charts either or both may be used as a [subchart](https://docs.helm.sh/chart_template_guide/#creating-a-subchart).
>
> Similarly, this Ingress chart can be combined -- either as a parent chart or a subchart -- with charts that define the services accepting traffic.

## Chart Details

The chart creates an Ingress resource for use with the [k8s-bigip-ctlr](http://clouddocs.f5.com/containers/latest/kubernetes/index.html).
