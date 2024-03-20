using Microsoft.AspNetCore.SignalR;

namespace Room;

public class Hub(RoomSystem roomSystem) : Microsoft.AspNetCore.SignalR.Hub
{
    private readonly RoomSystem roomSystem = roomSystem;

    public async Task CheckRoom()
    {
        try
        {
            // 当前用户获得房间信息
            var room = roomSystem.GetUser(Context.ConnectionId).Room;
            if (room != null) await Clients.Caller.SendAsync("OnMessage", room.ToModel().Serialize());
        }
        catch (Exception e) { Console.WriteLine(e); }
    }

    public async Task JoinRoom(int userId)
    {
        try
        {
            // 获得user对象 (已经由RPC加入房间)
            var user = roomSystem.Connect(Context.ConnectionId, userId);
            // 加入组
            string roomId = user.Room!.Id.ToString();
            await Groups.AddToGroupAsync(Context.ConnectionId, roomId);
            // 当前用户获得房间信息
            await Clients.Caller.SendAsync("OnMessage", user.Room!.ToModel().Serialize());
            // 通知房间其他人
            var msg = new Model.AddUser { user = user.ToModel() }.Serialize();
            await Clients.OthersInGroup(roomId).SendAsync("OnMessage", msg);
        }
        catch (Exception e) { Console.WriteLine(e); }
    }

    public async Task ExitRoom()
    {
        // 获得user对象
        var user = roomSystem.GetUser(Context.ConnectionId);
        // 退出房间
        var room = user.Room;
        if (room is null) return;
        room.RemoveUser(user.Id);
        // 通知房间所有人
        string roomId = room.Id.ToString();
        var msg = new Model.RemoveUser { userId = user.Id, ownerId = room.Owner?.Id ?? 0 }.Serialize();
        await Clients.Group(roomId).SendAsync("OnMessage", msg);
        // 退出组
        await Groups.RemoveFromGroupAsync(Context.ConnectionId, roomId);
    }

    public async Task SetAlready(bool value)
    {
        var user = roomSystem.GetUser(Context.ConnectionId);
        user.Already = value;
        var msg = new Model.SetAlready { userId = user.Id, value = value }.Serialize();
        await Clients.Group(user.Room!.Id.ToString()).SendAsync("OnMessage", msg);
    }

    public async Task StartGame()
    {
        var room = roomSystem.GetUser(Context.ConnectionId).Room!;
        var msg = new Model.StartGame().Serialize();
        await Clients.Group(room.Id.ToString()).SendAsync("OnMessage", msg);
        room.StartGame();
    }

    public Task SendGameMessage(string message)
    {
        var user = roomSystem.GetUser(Context.ConnectionId);
        user.Room?.SendGameMessage(message);
        return Task.CompletedTask;
    }

    // private async void OnGameMessage(string Message)
    // {
    //     await Clients.Groups(room.Id.ToString()).SendAsync("OnMessage", x.Serialize())
    // }

    public override async Task OnDisconnectedAsync(Exception? exception)
    {
        Console.WriteLine("disconnect");
        await ExitRoom();
        roomSystem.DisConnect(Context.ConnectionId);
    }
}