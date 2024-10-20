# Kubernets Cluster init

The kubernetes cluster is initalized via ansible. The ansibe deployment is configured for the default deployment
for the terraform resources. Make sure to have an active vpn connection to the cluster.

## Init the cluster

To be able to run the playbooks create a virtual environment via

```bash
python3 -m venv .venv
```

Activate the virtual environment using

```bash
source .venv/bin/activate
```

and install the requirements

```bash
pip install -r requirements.txt
ansible-galaxy install -r requirments.yml
```

To initialize the cluster run

```bash
ansible-playbook -i inventory main.yml
```
