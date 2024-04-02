import { ApolloError } from "@apollo/client";
import { Alert } from "@mui/material";
interface Props {
  error?: ApolloError;
}
export default function AlertApolloError({ error }: Props) {
  if (!error) {
    return <></>;
  }
  return (
    <Alert variant="outlined" severity="error">
      {error?.graphQLErrors.map((err) => err.message).join(",")}{" "}
    </Alert>
  );
}
