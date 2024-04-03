"use client";

import { GRAPHQL } from "@/config-global";
import { ApolloLink, HttpLink } from "@apollo/client";
import { setContext } from "@apollo/client/link/context";
import {
  ApolloNextAppProvider,
  NextSSRApolloClient,
  NextSSRInMemoryCache,
  SSRMultipartLink,
} from "@apollo/experimental-nextjs-app-support/ssr";
import { Authentication } from ".";
import { getStorage } from "@/utils/storage-available";

export const authLink = setContext((operation, { headers }) => {
  const auth = getStorage("auth") as
    | (Authentication & { projectId: string })
    | null;
  return {
    headers: {
      ...headers,
      ...(auth?.token && { authorization: `Bearer ${auth.token}` }),
      ...(operation.operationName !== "signin" &&
        auth?.projectId && { "x-project-id": auth.projectId }),
    },
  };
});
export const httpLink = new HttpLink({
  uri: GRAPHQL.endpoint,
  credentials: "omit",
});
export const apolloCache = new NextSSRInMemoryCache();
function makeClient() {
  return new NextSSRApolloClient({
    cache: apolloCache,
    credentials: "include",
    link:
      typeof window === "undefined"
        ? ApolloLink.from([
            new SSRMultipartLink({
              stripDefer: true,
            }),
            authLink,
            httpLink,
          ])
        : ApolloLink.from([authLink, httpLink]),
  });
}

export function ApolloWrapper({ children }: React.PropsWithChildren) {
  return (
    <ApolloNextAppProvider makeClient={makeClient}>
      {children}
    </ApolloNextAppProvider>
  );
}
