using Grpc.Core;
using GrpcRoomServer;
using Room;

namespace Room.Services;

public class RoomGrpcService(RoomSystem roomSystem) : RoomService.RoomServiceBase
{
    private readonly RoomSystem roomSystem = roomSystem;

    // public override Task<HelloReply> SayHello(HelloRequest request, ServerCallContext context)
    // {
    //     return Task.FromResult(new HelloReply
    //     {
    //         Message = "Hello " + request.Name
    //     });
    // }
    private const string url = "http://192.168.1.5:5001/hub";
    public override Task<GrpcRoomServer.Room> CreateRoom(CreateRoomRequest request, ServerCallContext context)
    {
        var room = roomSystem.CreateRoom(request.UserId);
        // Console.WriteLine("create " + request.UserId);
        return Task.FromResult(new GrpcRoomServer.Room { Id = room.Id, Url = url });
    }

    public override Task<GrpcRoomServer.Room> JoinRoom(JoinRoomRequest request, ServerCallContext context)
    {
        var room = roomSystem.AutoJoinRoom(request.UserId);
        return Task.FromResult(new GrpcRoomServer.Room { Id = room.Id, Url = url });
    }
}
