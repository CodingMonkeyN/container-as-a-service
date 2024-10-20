using ContainerAsAService.Application.Pods.Queries.GetPodLogs;
using ContainerAsAService.Application.Pods.Queries.GetPods;

namespace ContainerAsAService.Web.Endpoints;

public class Pods : EndpointGroupBase
{
    public override void Map(WebApplication app)
    {
        app.MapGroup(this)
            .AllowAnonymous()
            .MapGet(GetPods)
            .MapGet(GetPodLogs, "logs/{podNamespace}/{podName}");
    }

    private async Task<IResult> GetPods(ISender sender, [AsParameters] GetPodsQuery query)
    {
        return Results.Ok(await sender.Send(query));
    }

    private async Task<IResult> GetPodLogs(ISender sender, string podNamespace, string podName)
    {
        return Results.Stream(
            await sender.Send(new GetPodLogsQuery { PodNamespace = podNamespace, PodName = podName }));
    }
}
