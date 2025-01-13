using Api.Contracts;
using Api.Models;
using Api.Utils;
using k8s;

var builder = WebApplication.CreateBuilder(args);

var config = KubernetesClientConfiguration.BuildConfigFromConfigFile();
IKubernetes client = new Kubernetes(config);
builder.Services.AddOpenApi();
builder.Services.AddCors();
builder.Services.AddSingleton(client);
var app = builder.Build();

if (app.Environment.IsDevelopment()) app.MapOpenApi();

app.UseHttpsRedirection();
app.UseCors(corsBuilder => corsBuilder.AllowAnyOrigin().AllowAnyMethod().AllowAnyHeader());

app.MapGet("/pods", async (IKubernetes kubeClient) =>
    {
        var pods = await kubeClient.CoreV1.ListPodForAllNamespacesAsync();
        if (pods is null) return Results.NotFound();

        if (pods.Items is null) return Results.NotFound();
        return Results.Ok(pods.Items.Where(item => item.Spec.RuntimeClassName == "kata-qemu").Select(item =>
        {
            return new PodResponse
            {
                Name = item.Metadata.Name,
                Namespace = item.Metadata.NamespaceProperty,
                Cpu = item.Spec.Containers[0].Resources?.Limits?["cpu"].Value ?? "",
                Memory = item.Spec.Containers[0].Resources?.Limits?["memory"].Value ?? "",
                Port = item.Spec.Containers[0]?.Ports?[0].ContainerPort == null
                    ? 0
                    : item.Spec.Containers[0].Ports[0].ContainerPort,
                Ready = item.Status.ContainerStatuses[0].Ready,
                Status = Utils.MapStatus(item.Status.ContainerStatuses[0].State)
            };
        }));
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
                newResource, "apps.com.coding-monkey", "v1", request.Namespace, "containerdeployments"
            );
            return Results.Created();
        }
        catch (Exception e)
        {
            return Results.BadRequest(e.Message);
        }
    });

app.MapPost("updateDeployment",
    async (IKubernetes kubeClient, CreateDeploymentRequest request) =>
    {
        var currentObject = await kubeClient.CustomObjects
            .GetNamespacedCustomObjectAsync<ContainerDeployment>("apps.com.coding-monkey", "v1", request.Namespace,
                "containerdeployments", request.Name);

        if (currentObject is null) return Results.BadRequest();


        var newResource = new
        {
            apiVersion = "apps.com.coding-monkey/v1",
            kind = "ContainerDeployment",
            metadata = new
            {
                resourceVersion = currentObject.Metadata.ResourceVersion, name = request.Name,
                namespaceProperty = request.Namespace
            },
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
            await kubeClient.CustomObjects.ReplaceNamespacedCustomObjectAsync(
                newResource, "apps.com.coding-monkey", "v1", request.Namespace, "containerdeployments", request.Name
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
    var containerDeploymentList = await kubeClient.CustomObjects.ListClusterCustomObjectAsync<ContainerDeploymentList>(
        "apps.com.coding-monkey", "v1",
        "containerdeployments");
    return Results.Ok(containerDeploymentList.Items.Select(item => new ContainerDeploymentResponse
    {
        Namespace = item.Metadata.Namespace,
        Name = item.Metadata.Name,
        Cpu = item.Spec.Cpu,
        Memory = item.Spec.Memory,
        Image = item.Spec.Image,
        Replicas = item.Spec.Replicas,
        Port = item.Spec.Port,
        EnvironmentVariables = item.Spec.Env
    }));
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
    return Results.Ok(namespaces.Items.Select(item => item.Metadata.Name));
});

app.Run();