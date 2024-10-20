using ContainerAsAService.Application.Namespace.Queries.GetClusterNamespaces;

namespace ContainerAsAService.Web.Endpoints;

public class Namespace : EndpointGroupBase
{
    public override void Map(WebApplication app)
    {
        app.MapGroup(this)
            .AllowAnonymous()
            .MapGet(GetClusterNamespaces);
    }

    private async Task<IResult> GetClusterNamespaces(ISender sender, [AsParameters] GetClusterNamespacesQuery query)
    {
        return Results.Ok(await sender.Send(query));
    }
}
