---
description: This rule ensures proper access control for user data operations.
  It should be applied whenever implementing user data modification endpoints.
alwaysApply: true
---

Always implement access control checks that ensure users can only modify their own data unless they are administrators. Use middleware to validate that either: 1) The authenticated user is modifying their own data, OR 2) The authenticated user is an administrator