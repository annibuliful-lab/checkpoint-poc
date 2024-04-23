"use client";

import { GRAPHQL } from "@/config-global";
import {
  ApolloLink,
  FetchResult,
  Observable,
  ServerError,
} from "@apollo/client";
import { setContext } from "@apollo/client/link/context";
import {
  ApolloNextAppProvider,
  NextSSRApolloClient,
  NextSSRInMemoryCache,
  SSRMultipartLink,
} from "@apollo/experimental-nextjs-app-support/ssr";
import { Authentication } from ".";
import { getStorage } from "@/utils/storage-available";
import createUploadLink from "apollo-upload-client/createUploadLink.mjs";
import { onError } from "@apollo/client/link/error";
import { GraphQLError } from "graphql";

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

const callRefreshToken = async () => {
  try {
    const auth = getStorage("auth") as
      | (Authentication & { projectId: string })
      | null;

    const myHeaders = new Headers();
    if (auth?.projectId) {
      myHeaders.append("x-project-id", auth?.projectId);
    }

    const refreshTokenResponse = await (
      await fetch(GRAPHQL.endpoint, {
        method: "POST",
        headers: myHeaders,
        body: JSON.stringify({
          query:
            "mutation callRefreshToken($refreshToken: String!) {\n  refreshToken(refreshToken: $refreshToken) {\n    token\n    refreshToken\n    userId\n  }\n}\n",
          variables: {
            refreshToken: auth?.refreshToken,
          },
        }),
        redirect: "follow",
      })
    ).json();

    const accessToken = refreshTokenResponse.data?.refreshToken?.token;
    if (!accessToken) {
      throw "not found access token.";
    }
    
    window.localStorage.setItem(
      "auth",
      JSON.stringify({
        ...refreshTokenResponse?.data?.refreshToken,
        projectId: auth?.projectId,
      })
    );
    return accessToken;
  } catch (err) {
    window.localStorage.clear();
    window.location.reload();
    throw err;
  }
};

const errorLink = onError((response) => {
  const { networkError, operation, forward } = response;
  const networkErrors = networkError as ServerError;
  if (networkErrors?.statusCode === 401) {
    if (operation.operationName === "refreshToken") return;
    const observable = new Observable<FetchResult<Record<string, unknown>>>(
      (observer) => {
        (async () => {
          try {
            const accessToken = await callRefreshToken();

            if (!accessToken) {
              throw new GraphQLError("Invalid or expire token!");
            }

            const subscriber = {
              next: observer.next.bind(observer),
              error: observer.error.bind(observer),
              complete: observer.complete.bind(observer),
            };

            forward(operation).subscribe(subscriber);
          } catch (err) {
            observer.error(err);
          }
        })();
      }
    );

    return observable;
  }
});

export const httpLink = createUploadLink({
  uri: GRAPHQL.endpoint,
  credentials: "omit",
}) as unknown as ApolloLink;
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
            errorLink,
            authLink,
            httpLink,
          ])
        : ApolloLink.from([errorLink, authLink, httpLink]),
  });
}

export function ApolloWrapper({ children }: React.PropsWithChildren) {
  return (
    <ApolloNextAppProvider makeClient={makeClient}>
      {children}
    </ApolloNextAppProvider>
  );
}
