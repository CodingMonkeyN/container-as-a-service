using ContainerAsAService.Application.Common.Interfaces;

namespace ContainerAsAService.Application.Namespace.Queries.GetClusterNamespaces;

public record GetClusterNamespacesQuery : IRequest<IList<string>>;

public class GetClusterNamespacesQueryHandler(IKubernetesService kubernetesService)
    : IRequestHandler<GetClusterNamespacesQuery, IList<string>>
{
    public async Task<IList<string>> Handle(GetClusterNamespacesQuery request, CancellationToken cancellationToken)
    {
        return await kubernetesService.GetClusterNamespacesAsync();
    }
}
