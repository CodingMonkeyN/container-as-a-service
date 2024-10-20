using ContainerAsAService.Application.Common.Interfaces;

namespace ContainerAsAService.Application.Pods.Queries.GetPodLogs;

public record GetPodLogsQuery : IRequest<Stream>
{
    public required string PodNamespace { get; init; }
    public required string PodName { get; init; }
}

public class GetPodLogsQueryValidator : AbstractValidator<GetPodLogsQuery>
{
    public GetPodLogsQueryValidator()
    {
        RuleFor(x => x.PodNamespace).NotEmpty();
        RuleFor(x => x.PodName).NotEmpty();
    }
}

public class GetPodLogsQueryHandler(IKubernetesService kubernetesService) : IRequestHandler<GetPodLogsQuery, Stream>
{
    public async Task<Stream> Handle(GetPodLogsQuery request, CancellationToken cancellationToken)
    {
        return await kubernetesService.GetPodLogsAsync(request.PodNamespace, request.PodName);
    }
}
