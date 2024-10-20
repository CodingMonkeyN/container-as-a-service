using ContainerAsAService.Application.ContainerDeployment.Queries.GetContainerDeployments;
using ContainerAsAService.Application.Models;

namespace ContainerAsAService.Application.Common.Interfaces;

public interface IContainerDeploymentService
{
    Task<IList<ContainerDeploymentDto>> GetContainerDeploymentsAsync();

    Task CreateContainerDeploymentAsync(CreateContainerDeployment createContainerDeployment);

    Task UpdateContainerDeploymentAsync(UpdateContainerDeployment updateContainerDeployment);
}
