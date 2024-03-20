using System.Reflection;
using Microsoft.AspNetCore.Server.Kestrel.Core;
using Room;
using Room.Services;

using var srg = new StreamReader("Static/general.json");
Model.General.Init(srg.ReadToEnd());
using var src = new StreamReader("Static/card.json");
Model.Card.Init(src.ReadToEnd());

Assembly.Load("Skills");

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
// 房间系统
builder.Services.AddSingleton<RoomSystem>();
// Grpc
builder.Services.AddGrpc();
// builder.Services.AddControllers();
// JWT配置
// builder.Services.Configure<JWTOptions>(builder.Configuration.GetSection("JWT"));
// // 
// builder.Services.AddAuthentication(JwtBearerDefaults.AuthenticationScheme).AddJwtBearer(x =>
// {
//     var jwtOpt = builder.Configuration.GetSection("JWT").Get<JWTOptions>();
//    var a= builder.Configuration["JWT:Secret"];
// #pragma warning disable CS8604 // 引用类型参数可能为 null。
//     byte[] keyBytes = Encoding.UTF8.GetBytes(a);
// #pragma warning restore CS8604 // 引用类型参数可能为 null。
//     var secKey = new SymmetricSecurityKey(keyBytes);
//     x.TokenValidationParameters = new()
//     {
//         ValidateIssuer = false,
//         ValidateAudience = false,
//         ValidateLifetime = true,
//         ValidateIssuerSigningKey = true,
//         IssuerSigningKey = secKey
//     };
//     x.Events = new JwtBearerEvents
//     {
//         OnMessageReceived = context =>
//         {
//             var accessToken = context.Request.Query["access_token"];
//             var path = context.HttpContext.Request.Path;
//             if (!string.IsNullOrEmpty(accessToken) &&
//                 (path.StartsWithSegments("/Hubs/ChatRoomHub")))
//             {
//                 context.Token = accessToken;
//             }
//             return Task.CompletedTask;
//         }
//     };
// });
builder.Services.AddSignalR();

builder.WebHost.ConfigureKestrel((options) =>
{
    options.ListenAnyIP(5001); // 设置端口为 5000
    options.ListenAnyIP(5002, listenOptions =>
    {
        // 设置 gRPC 的端口为 5001，并启用 HTTP/2 协议
        listenOptions.Protocols = HttpProtocols.Http2;
    });
});

var app = builder.Build();

// Configure the HTTP request pipeline.
app.MapHub<Hub>("/hub");
app.MapGrpcService<RoomGrpcService>();

// app.UseEndpoints().useurl

// app.MapControllers();

app.Run();
