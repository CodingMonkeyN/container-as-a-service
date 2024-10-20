namespace ContainerAsAService.Infrastructure.Models;

public class ContainerDeploymentList
{
    public required string ApiVersion { get; init; }
    public required List<ContainerDeployment> Items { get; init; }
}

public class ContainerDeployment
{
    public required string ApiVersion { get; init; }
    public required ContainerMetadata Metadata { get; init; }
    public required ContainerDeploymentSpec Spec { get; init; }
}

public class ContainerDeploymentSpec
{
    public required string Cpu { get; init; }
    public required string Memory { get; init; }
    public required string Image { get; init; }
    public required int Replicas { get; init; }
    public required int Port { get; init; }
    public Dictionary<string, string>? Env { get; init; }
}

public class ContainerMetadata
{
    public required string Name { get; init; }
    public required string Namespace { get; init; }
    public required string ResourceVersion { get; init; }
}
