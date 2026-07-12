# Generated OpenAPI types (optional)

Hand-maintained types live in `src/api/types.ts`.

To regenerate from the backend Swagger spec:

```bash
npm run generate:api
```

This reads `../iam-backend/docs/swagger.yaml` and writes `schema.ts` here.
