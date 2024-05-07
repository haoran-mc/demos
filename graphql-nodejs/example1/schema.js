import { makeExecutableSchema } from '@graphql-tools/schema'
 
export const schema = makeExecutableSchema({
  typeDefs: /* GraphQL */ `
    type Query {
      hello: String
    }
  `,
  resolvers: {
    Query: {
      hello: () => 'world'
    }
  }
})