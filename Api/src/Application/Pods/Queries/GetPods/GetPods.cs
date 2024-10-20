using ContainerAsAService.Application.Common.Interfaces;

namespace ContainerAsAService.Application.Pods.Queries.GetPods;

public record GetPodsQuery : IRequest<IList<PodDto>>;

public class GetPodsQueryHandler(IKubernetesService kubernetesService)
    : IRequestHandler<GetPodsQuery, IList<PodDto>>
{
    public async Task<IList<PodDto>> Handle(GetPodsQuery request, CancellationToken cancellationToken)
    {
        return await kubernetesService.GetPodsAsync();
    }
}
