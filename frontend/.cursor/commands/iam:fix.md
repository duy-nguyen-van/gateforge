Fix the described bug or issue systematically.

**Rule: Find root cause FIRST. No random fixes.**

## Workflow

### Step 1: Understand
- Read error, stack trace, or reproduction steps
- Check browser Network/Application tabs for auth issues

### Step 2: Root Cause Analysis
- Trace route guard → auth provider → api client → backend
- `git log --oneline -10` for recent changes
- 2–3 hypotheses; verify with code/network/console

### Step 3: Implement Fix
- Minimal change at correct layer
- Match existing `ApiError` and token patterns

### Step 4: Verify
- `make lint` and `make build`
- Manual reproduction of fixed flow

### Step 5: Self-Review
- Right layer? Token storage safe? Guards correct?

## Red Flags
- "Just try changing this" — stop, find root cause
- "Build passes" — confirm the actual bug is fixed in browser

Suggest conventional commit message when done. Do not commit unless asked.
