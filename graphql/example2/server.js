// --> 1. GraphQL
const typeDefs = `
  type Author {
    id: ID! # the ! means that every author object _must_ have an id
    firstName: String
    lastName: String
    """
    the list of Posts by this author
    """
    posts: [Post]
  }

  type Post {
    id: ID!
    title: String
    author: Author
    votes: Int
  }

  # the schema allows the following query:
  type Query {
    posts: [Post]
    authors: [Author]
    post(id: ID!): Post
    author(id: ID!): Author
  }  

  # this schema allows the following mutation:
  type Mutation {
    upvotePost(postId: ID!): Post
  }

  # we need to tell the server which types represent the root query
  # and root mutation types. We call them RootQuery and RootMutation by convention.
  schema {
    query: Query
    mutation: Mutation
  }
`

const authors = [
  { id: '1', firstName: 'John', lastName: 'Doe' },
  { id: '2', firstName: 'Jane', lastName: 'Smith' },
];

const posts = [
  { id: '1', title: 'GraphQL Basics', authorId: '1', votes: 10 },
  { id: '2', title: 'Advanced GraphQL Techniques', authorId: '1', votes: 5 },
  { id: '3', title: 'Introduction to React', authorId: '2', votes: 15 },
  { id: '4', title: 'State Management in React', authorId: '2', votes: 8 },
];

// import { find, filter } from 'lodash'; // 添加导入

import lodash from 'lodash';
const { find, filter } = lodash;

// --> 2. resolvers
const resolvers = {
  Query: {
    posts() { // 添加查询贴子的函数
      return posts
    },
    authors() { // 添加查询作者的函数
      return authors
    },
    post(_, { id }) { // 查询特定贴子的函数
      return find(posts, { id });
    },
    author(_, { id }) { // 查询特定作者的函数
      return find(authors, { id });
    }
  },
  Mutation: {
    upvotePost(_, { postId }) {
      const post = find(posts, { id: postId })
      if (!post) {
        throw new Error(`Couldn't find post with id ${postId}`)
      }
      post.votes += 1
      return post
    }
  },
  // 作者
  Author: {
    // 作者写过的贴子
    posts(author) {
      return filter(posts, { authorId: author.id })
    }
  },
  // 贴子
  Post: {
    // 通过贴子，得到作者
    author(post) {
      return find(authors, { id: post.authorId })
    }
  }
}

// --> 3. the schema and resolvers are combined using makeExecutableSchema
import { makeExecutableSchema } from '@graphql-tools/schema'

const executableSchema = makeExecutableSchema({
  typeDefs,
  resolvers
})

// --> 4. GraphQL-Tools schema can be consumed by frameworks like GraphQL Yoga, Apollo GraphQL or express-graphql
import { createServer } from 'node:http'
import { createYoga } from 'graphql-yoga'

const yoga = createYoga({
  schema: executableSchema
})

const server = createServer(yoga)

server.listen(4000, () => {
  console.log('Yoga is listening at http://localhost:4000/graphql')
})