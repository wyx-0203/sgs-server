using Microsoft.AspNetCore.SignalR;
using Model;

namespace Room;

public class RoomSystem(IHubContext<Hub> hubContext)
{
    private readonly IHubContext<Hub> hubContext = hubContext;

    public List<Room> Rooms { get; } = [];
    public List<User> Users { get; } = [];
    private int index = 1;

    public Room CreateRoom(int userId)
    {
        var room = new Room { Id = index++, Mode = Mode._3V3, HubContext = hubContext };
        Rooms.Add(room);
        var user = GetUser(userId);
        room.AddUser(user);
        return room;
    }

    public Room AutoJoinRoom(int userId)
    {
        var room = Rooms.Find(x => !x.Full);
        if (room is null) return CreateRoom(userId);
        // room??=CreateRoom();
        var user = GetUser(userId);
        room.AddUser(user);

        // 通知房间其他人
        // hubContext.Clients.Group(room.Id.ToString()).SendAsync("OnMessage", room.ToModel().Serialize());

        return room;
    }

    private readonly Dictionary<string, User> userMap = [];

    public User Connect(string connectionId, int userId)
    {
        if (userMap.TryGetValue(connectionId, out User? user)) return user;

        // Console.WriteLine(userId);
        user = Rooms.SelectMany(x => x.Users).First(x => x?.Id == userId)!;
        userMap.Add(connectionId, user);
        return user;
    }

    public void DisConnect(string connectionId) => userMap.Remove(connectionId);

    public User GetUser(string connectionId) => userMap[connectionId];

    private User GetUser(int userId)
    {
        var user = Users.Find(x => x.Id == userId);
        if (user != null) return user;
        user = new User { Id = userId };
        Users.Add(user);
        return user;
    }
    // {
    //     return userMap.TryGetValue(connectionId, out User? value) ? value : null;
    // }
}

public class Room
{
    public required int Id { get; set; }
    public required Mode Mode { get; set; }
    public User?[] Users { get; } = new User[2];
    public User? Owner { get; private set; }
    // public int count=0;
    public bool Full => Users.All(x => x != null);

    public GameCore.Game? Game { get; private set; }
    public required IHubContext<Hub> HubContext { get; set; }

    public void AddUser(User user)
    {
        if (Full) return;
        for (int i = 0; i < Users.Length; i++)
        {
            if (Users[i] != null) continue;
            
            Users[i] = user;
            user.Room = this;
            user.Position = i;
            Owner ??= user;
            Console.WriteLine($"adduser roomId:{Id}, userId:{user.Id}");
            break;
        }
    }

    public void RemoveUser(int userId)
    {
        for (int i = 0; i < Users.Length; i++)
        {
            if (Users[i]?.Id == userId)
            {
                var user = Users[i]!;
                Users[i] = null;
                user.Room = null;
                if (Owner == user) Owner = Users.FirstOrDefault(x => x != null && x != user);
                Console.WriteLine($"exitroom roomId:{Id}, userId:{user.Id}, ownerId:{Owner?.Id}");
                return;
            }
        }
    }

    public async void StartGame()
    {
        if (!Users.All(x => x != null)) return;
        var ids = Users.Select(x => x!.Id).Shuffle();
        Game = new GameCore.Game(Mode, ids, OnSendToClient);
        
        await Game.Init();
        await Game.Run();

        foreach (var i in Users) if (i != null) i.Already = false;
    }

    private async void OnSendToClient(Message message)
    {
        Console.WriteLine(message.Serialize());
        await HubContext.Clients.Group(Id.ToString()).SendAsync("OnMessage", message.Serialize());
    }

    public void SendGameMessage(string message)
    {
        if (message.DeSerialize() is Message message1) Game?.eventSystem.PushMessage(message1);
    }

    public JoinRoom ToModel() => new()
    {
        mode = "3v3",
        users = Users.Where(x => x != null).Select(x => x!.ToModel()).ToList(),
        ownerId = Owner!.Id
    };
}

public class User
{
    public string? ConnectionId { get; set; }
    public required int Id { get; set; }
    public Room? Room { get; set; }
    public int Position { get; set; }
    public bool Already { get; set; }

    public Model.User ToModel() => new()
    {
        id = Id,
        position = Position,
        already = Already
    };
}