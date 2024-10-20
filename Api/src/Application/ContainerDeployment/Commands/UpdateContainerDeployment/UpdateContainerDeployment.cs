using ContainerAsAService.Application.Common.Interfaces;

namespace ContainerAsAService.Application.ContainerDeployment.Commands.UpdateContainerDeployment;

public record UpdateContainerDeploymentCommand : IRequest
{
    public required string Name { get; init; }
    public required string Namespace { get; init; }
    public required string Image { get; init; }
    public required string Cpu { get; init; }
    public required string Memory { get; init; }
    public required int Replicas { get; init; }
    public required int Port { get; init; }
    public Dictionary<string, string>? EnvironmentVariables { get; init; }
}

public class UpdateContainerDeploymentCommandValidator : AbstractValidator<UpdateContainerDeploymentCommand>
{
    public UpdateContainerDeploymentCommandValidator()
    {
        RuleFor(x => x.Namespace).NotEmpty();
        RuleFor(x => x.Image).NotEmpty();
        RuleFor(x => x.Name).NotEmpty().MaximumLength(15);
        RuleFor(x => x.Cpu).NotEmpty();
        RuleFor(x => x.Memory).NotEmpty();
        RuleFor(x => x.Replicas).NotEmpty();
        RuleFor(x => x.Port).NotEmpty();
    }
}

public class UpdateContainerDeploymentCommandHandler(IContainerDeploymentService containerDeploymentService)
    : IRequestHandler<UpdateContainerDeploymentCommand>
{
    public async Task Handle(UpdateContainerDeploymentCommand request, CancellationToken cancellationToken)
    {
        await containerDeploymentService.UpdateContainerDeploymentAsync(new Models.UpdateContainerDeployment
        {
            Namespace = request.Namespace,
            Image = request.Image,
            Cpu = request.Cpu,
            EnvironmentVariables = request.EnvironmentVariables,
            Memory = request.Memory,
            Name = request.Name,
            Port = request.Port,
            Replicas = request.Replicas
        });
    }
}
