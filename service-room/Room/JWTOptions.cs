public class JWTOptions
{
    public required string Secret { get; set; }
    public required string ExpireSeconds { get; set; }
}