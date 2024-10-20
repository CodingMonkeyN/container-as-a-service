using ContainerAsAService.Application.Common.Interfaces;

namespace ContainerAsAService.Application.ContainerDeployment.Commands.CreateContainerDeployment;

public record CreateContainerDeploymentCommand : IRequest
{
    public required string Namespace { get; init; }
    public required string Name { get; init; }
    public required string Cpu { get; init; }
    public required string Memory { get; init; }
    public required string Image { get; init; }
    public required int Replicas { get; init; }
    public required int Port { get; init; }
    public Dictionary<string, string>? EnvironmentVariables { get; init; }
}

public class CreateContainerDeploymentCommandValidator : AbstractValidator<CreateContainerDeploymentCommand>
{
    public CreateContainerDeploymentCommandValidator()
    {
        RuleFor(x => x.Namespace).NotEmpty();
        RuleFor(x => x.Name).NotEmpty().MaximumLength(15);
        RuleFor(x => x.Cpu).NotEmpty();
        RuleFor(x => x.Memory).NotEmpty();
        RuleFor(x => x.Image).NotEmpty();
        RuleFor(x => x.Replicas).NotEmpty();
        RuleFor(x => x.Port).NotEmpty();
    }
}

public class CreateContainerDeploymentCommandHandler(IContainerDeploymentService containerDeploymentService)
    : IRequestHandler<CreateContainerDeploymentCommand>
{
    public async Task Handle(CreateContainerDeploymentCommand request, CancellationToken cancellationToken)
    {
        await containerDeploymentService.CreateContainerDeploymentAsync(new Models.CreateContainerDeployment
        {
            Cpu = request.Cpu,
            EnvironmentVariables = request.EnvironmentVariables,
            Image = request.Image,
            Memory = request.Memory,
            Name = request.Name,
            Namespace = request.Namespace,
            Port = request.Port,
            Replicas = request.Replicas
        });
    }
}
