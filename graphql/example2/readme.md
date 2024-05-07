官方提供的运行不起来的样例，进行了完善。

schema 和 server 代码在一个文件里。


- 查询所有帖子和它们的作者：
```
query {
  posts {
    id
    title
    author {
      id
      firstName
      lastName
    }
    votes
  }
}
```

- 查询所有作者及其帖子：
```
query {
  authors {
    id
    firstName
    lastName
    posts {
      id
      title
      votes
    }
  }
}
```

- 查询指定帖子的信息：
```
query {
  post(id: "1") {
    id
    title
    author {
      id
      firstName
      lastName
    }
    votes
  }
}
```

- 为帖子点赞：
```
mutation {
  upvotePost(postId: "1") {
    id
    title
    votes
  }
}
```