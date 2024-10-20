using ContainerAsAService.Application.ContainerDeployment.Commands.CreateContainerDeployment;
using ContainerAsAService.Application.ContainerDeployment.Commands.UpdateContainerDeployment;
using ContainerAsAService.Application.ContainerDeployment.Queries.GetContainerDeployments;

namespace ContainerAsAService.Web.Endpoints;

public class ContainerDeployment : EndpointGroupBase
{
    public override void Map(WebApplication app)
    {
        app.MapGroup(this)
            .AllowAnonymous()
            .MapPost(CreateContainerDeployment, "create")
            .MapPost(UpdateContainerDeployment, "update")
            .MapGet(GetContainerDeployments);
    }

    private async Task<IResult> GetContainerDeployments(ISender sender,
        [AsParameters] GetContainerDeploymentQuery query)
    {
        return Results.Ok(await sender.Send(query));
    }

    private async Task<IResult> CreateContainerDeployment(ISender sender, CreateContainerDeploymentCommand command)
    {
        await sender.Send(command);
        return Results.NoContent();
    }

    private async Task<IResult> UpdateContainerDeployment(ISender sender, UpdateContainerDeploymentCommand command)
    {
        await sender.Send(command);
        return Results.NoContent();
    }
}
