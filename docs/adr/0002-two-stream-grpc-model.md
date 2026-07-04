# Two-stream gRPC model for real-time message delivery

Separate real-time delivery into two distinct gRPC server-streaming RPCs: one `ActiveChannelStream` that pushes full message events for the single channel the user is currently focused on, and one `NotificationStream` that pushes lightweight "new message" pings (channel_id + author + preview snippet) for all other channels the user can see.

**Why two streams instead of one?** A single multi-plexed stream would force the TUI client to filter out full message payloads for unfocused channels (wasted bandwidth and parsing). Per-channel streams would require the client to open and close N separate gRPC connections when switching focus. The two-stream split means the TUI opens exactly two long-lived streams per session: one gets rich payloads for the active view, the other gets just enough metadata to render a notification badge in the channel list.

**Status**: accepted

**Considered Options**:
- **Single multi-plexed stream** — one stream carries all events tagged with channel_id. Client filters on receive. Simpler connection management but wastes bandwidth on full message payloads for non-focused channels.
- **Per-channel stream** — open/close a stream per channel as the user navigates. Cleanest separation but adds connection churn and per-stream overhead on the server.
- **Two-stream model (chosen)** — active channel gets full payloads, background channels get lightweight notifications. Best bandwidth/profile trade-off for a TUI that renders one channel at a time with badges on the rest.
