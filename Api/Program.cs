using Api.Models;
using k8s;

var builder = WebApplication.CreateBuilder(args);

var config = KubernetesClientConfiguration.BuildConfigFromConfigFile(); // Kubeconfig wird geladen
IKubernetes client = new Kubernetes(config);
builder.Services.AddOpenApi();
builder.Services.AddSingleton(client);
var app = builder.Build();

if (app.Environment.IsDevelopment()) app.MapOpenApi();

app.UseHttpsRedirection();


app.MapGet("/pods", async (IKubernetes kubeClient) =>
    {
        var pods = await kubeClient.CoreV1.ListPodForAllNamespacesAsync();
        Console.WriteLine("Pods im Cluster:");
        foreach (var pod in pods.Items)
            Console.WriteLine($"Namespace: {pod.Metadata.NamespaceProperty}, Name: {pod.Metadata.Name}");
    }
);

app.MapPost("/createDeployment",
    async (CreateDeploymentRequest request) =>
    {
        var newResource = new
        {
            apiVersion = "apps.com.coding-monkey/v1",
            kind = "ContainerDeployment",
            metadata = new { name = request.Name, namespaceProperty = request.Namespace },
            spec = new
            {
                image = request.Image,
                memory = request.Memory,
                cpu = request.Cpu,
                replicas = request.Replicas,
                port = request.Port,
                env = request.EnvironmentVariables
            }
        };

        try
        {
            await client.CustomObjects.CreateNamespacedCustomObjectAsync(
                newResource, "apps.com.coding-monkey", "v1", "ingress", "containerdeployments"
            );
            return Results.Created();
        }
        catch (Exception e)
        {
            return Results.BadRequest(e.Message);
        }
    });

app.MapGet("/containerDeployments", async (IKubernetes kubeClient) =>
{
    var x = await kubeClient.CustomObjects.ListClusterCustomObjectAsync("apps.com.coding-monkey", "v1",
        "containerdeployments");
    return Results.Ok(x);
});

app.MapGet("/logs/{namespaceName}/{podName}", async (IKubernetes kubeClient, string namespaceName, string podName) =>
{
    var logResponse = await kubeClient.CoreV1.ReadNamespacedPodLogAsync(
        podName,
        namespaceName
    );
    return Results.Stream(logResponse);
});

app.MapGet("/namespaces", async (IKubernetes kubeClient) =>
{
    var namespaces = await kubeClient.CoreV1.ListNamespaceAsync();
    return Results.Ok(namespaces);
});

app.Run();