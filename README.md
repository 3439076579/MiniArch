# MiniArch
## Web
Web小框架giu，是在模仿gin的基础上实现的Web框架，包含了中间件，路由前缀匹配树以及上下文Context和路由组的实现。其中Context的Keys使用了读写锁来保证并发安全，但是路由中的匹配树是线程不安全的。 
## ORM
数据库ORM小框架，是在学习并模仿gorm的基础上实现的ORM框架，包括进行数据库的迁移，支持事务，组装SQL语句，实现对象结构与数据库表的映射，还支持钩子操作。
## Distributed Cache
简单版分布式缓存，在学习并模仿Gocache的基础上实现的分布式缓存，基于json的格式进行服务器之间的传输，实现了LRU算法以及一致性哈希环，并使用`singleflight`来避免缓存穿透和缓存击穿，与Gocache相比，少了对过期时间的支持，以及无法检测每组Group的情况
