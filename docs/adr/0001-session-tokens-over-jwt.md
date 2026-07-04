# Session tokens over JWT for gRPC authentication

Session tokens instead of JWTs for gRPC auth. On login the server generates a random token from `crypto/rand`, stores `SHA-256(token_hash) + user_id + expiry` in the `sessions` table, and returns the raw token to the client. Every gRPC call carries it in the `authorization: bearer` metadata header; the server hashes it and looks up the session.

**Why not JWT?** JWT's statelessness is largely wasted here — every RPC already hits Postgres for channel/message data, so skipping one DB lookup for auth doesn't move the needle. Session tokens let us revoke individual sessions with a simple `DELETE` (kick a user, log out a specific device, force re-auth) without a blocklist that would bring back the DB dependency anyway. They also avoid key rotation, signing algorithm choices, and the risk of a leaked signing key forging tokens.

**Status**: accepted

**Considered Options**:
- **JWT (HS256)** — stateless, but needs a blocklist for revocation. Same DB dependency but now with key management overhead.
- **JWT (RS256)** — adds asymmetric key management and rotation for an app with one auth server and one DB. Overkill.
- **mTLS** — correct for machine-to-machine, wrong for a chat app with ephemeral TUI sessions.
- **Session tokens (chosen)** — simple `crypto/rand` token, hashed before storage, easy to revoke, no key management.
