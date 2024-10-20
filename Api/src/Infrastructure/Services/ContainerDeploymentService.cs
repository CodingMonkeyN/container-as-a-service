using ContainerAsAService.Application.Common.Interfaces;
using ContainerAsAService.Application.ContainerDeployment.Queries.GetContainerDeployments;
using ContainerAsAService.Application.Models;
using ContainerAsAService.Infrastructure.Models;
using k8s;

namespace ContainerAsAService.Infrastructure.Services;

public class ContainerDeploymentService(IKubernetes kubernetes)
    : IContainerDeploymentService
{
    public async Task<IList<ContainerDeploymentDto>> GetContainerDeploymentsAsync()
    {
        ContainerDeploymentList? containerDeploymentList =
            await kubernetes.CustomObjects.ListClusterCustomObjectAsync<ContainerDeploymentList>(
                "apps.com.coding-monkey", "v1",
                "containerdeployments");
        return containerDeploymentList.Items.Select(item => new ContainerDeploymentDto
        {
            Namespace = item.Metadata.Namespace,
            Name = item.Metadata.Name,
            Cpu = item.Spec.Cpu,
            Memory = item.Spec.Memory,
            Image = item.Spec.Image,
            Replicas = item.Spec.Replicas,
            Port = item.Spec.Port,
            EnvironmentVariables = item.Spec.Env
        }).ToList();
    }

    public async Task CreateContainerDeploymentAsync(CreateContainerDeployment createContainerDeployment)
    {
        var newResource = new
        {
            apiVersion = "apps.com.coding-monkey/v1",
            kind = "ContainerDeployment",
            metadata =
                new { name = createContainerDeployment.Name, namespaceProperty = createContainerDeployment.Namespace },
            spec = new
            {
                image = createContainerDeployment.Image,
                memory = createContainerDeployment.Memory,
                cpu = createContainerDeployment.Cpu,
                replicas = createContainerDeployment.Replicas,
                port = createContainerDeployment.Port,
                env = createContainerDeployment.EnvironmentVariables
            }
        };

        await kubernetes.CustomObjects.CreateNamespacedCustomObjectAsync(
            newResource, "apps.com.coding-monkey", "v1", createContainerDeployment.Namespace, "containerdeployments"
        );
    }

    public async Task UpdateContainerDeploymentAsync(UpdateContainerDeployment updateContainerDeployment)
    {
        ContainerDeployment? currentObject = await kubernetes.CustomObjects
            .GetNamespacedCustomObjectAsync<ContainerDeployment>("apps.com.coding-monkey", "v1",
                updateContainerDeployment.Namespace,
                "containerdeployments", updateContainerDeployment.Name);

        Guard.Against.Null(currentObject, nameof(currentObject));

        var newResource = new
        {
            apiVersion = "apps.com.coding-monkey/v1",
            kind = "ContainerDeployment",
            metadata = new
            {
                resourceVersion = currentObject.Metadata.ResourceVersion,
                name = updateContainerDeployment.Name,
                namespaceProperty = updateContainerDeployment.Namespace
            },
            spec = new
            {
                image = updateContainerDeployment.Image,
                memory = updateContainerDeployment.Memory,
                cpu = updateContainerDeployment.Cpu,
                replicas = updateContainerDeployment.Replicas,
                port = updateContainerDeployment.Port,
                env = updateContainerDeployment.EnvironmentVariables
            }
        };

        await kubernetes.CustomObjects.ReplaceNamespacedCustomObjectAsync(
            newResource, "apps.com.coding-monkey", "v1", updateContainerDeployment.Namespace, "containerdeployments",
            updateContainerDeployment.Name
        );
    }
}
