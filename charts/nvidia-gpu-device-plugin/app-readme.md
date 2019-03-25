# NVIDIA GPU Device Plugin

The NVIDIA device plugin for Kubernetes is a Daemonset that allows you to automatically:  
- Expose the number of GPUs on each nodes of your cluster
- Keep track of the health of your GPUs
- Run GPU enabled containers in your Kubernetes cluster.

The app will use `nodeselector` to specify which node to deploy the device plugin.  
**For deploying the default NVIDIA device plugin, please set the following label to the cluster-node**  
` gpu.rancher.io/type="default" `

**For deploying the GPU-Sharing supported device plugin, please set the following label to the cluster-node**  
` gpu.rancher.io/type="shared" `

The device plugin should be deployed once.