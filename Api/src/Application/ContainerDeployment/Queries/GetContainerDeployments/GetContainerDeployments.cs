using ContainerAsAService.Application.Common.Interfaces;

namespace ContainerAsAService.Application.ContainerDeployment.Queries.GetContainerDeployments;

public record GetContainerDeploymentQuery : IRequest<IList<ContainerDeploymentDto>>;

public class GetContainerDeploymentQueryHandler(IContainerDeploymentService containerDeploymentService)
    : IRequestHandler<GetContainerDeploymentQuery, IList<ContainerDeploymentDto>>
{
    public async Task<IList<ContainerDeploymentDto>> Handle(GetContainerDeploymentQuery request,
        CancellationToken cancellationToken)
    {
        return await containerDeploymentService.GetContainerDeploymentsAsync();
    }
}
