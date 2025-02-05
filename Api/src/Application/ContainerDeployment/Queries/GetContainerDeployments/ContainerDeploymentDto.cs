namespace ContainerAsAService.Application.ContainerDeployment.Queries.GetContainerDeployments;

public record ContainerDeploymentDto
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
