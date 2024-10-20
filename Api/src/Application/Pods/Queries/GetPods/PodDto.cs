namespace ContainerAsAService.Application.Pods.Queries.GetPods;

public record PodDto
{
    public required string Namespace { get; init; }
    public required string Name { get; init; }
    public required string Cpu { get; init; }
    public required string Memory { get; init; }
    public required int Port { get; init; }
    public required string Status { get; init; }
    public required bool Ready { get; init; }
}
