import type { CodegenConfig } from "@graphql-codegen/cli";

const config: CodegenConfig = {
  overwrite: true,
  schema: "./generated.graphql",
  documents: "./dashboard/src/apollo-client/**/**/*.gql",
  generates: {
    "./dashboard/src/apollo-client/generated.ts": {
      // preset: 'client',
      plugins: [
        "typescript",
        "typescript-operations",
        "typescript-react-apollo",
      ],
      config: {
        withHooks: true,
        withRefetchFn: true,
        withMutationFn: true,
      },
    },
  },
};

export default config;
