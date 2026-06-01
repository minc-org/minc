package providers

import "time"

// MicroShiftServiceMaxRetries bounds how long we poll systemctl is-active for the microshift
// unit inside the container. While the unit is "activating", is-active exits non-zero (often 3),
// so each poll is a retry. Rootless or slow storage can keep the unit activating for many minutes.
//
// pkg/retry uses linear backoff: sleep = InitialDelay * attempt.
// Total wait = InitialDelay * (1+2+...+(N-1)) = 2s * (17*18/2) = 306s ≈ 5 min.
const MicroShiftServiceMaxRetries = 18

const MicroShiftServiceInitialRetryDelay = 2 * time.Second
