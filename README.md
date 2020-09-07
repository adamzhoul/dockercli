# Background

Docker CLI or kubectl are both great tools. But:

	1. they are designed for admin, not the developer. 
	2. developers know less about docker or k8s.
	3. developers have no permission to run kubectl command, especially on product env.
	4. kubectl uses `kubelet`  to receive and send data, which adds its heavy.

So, we need a tool for developers.
And what developers need are:

	1. login to docker to see code,env...
	2. check the log using `tail` or grep .. instead of kibana.
    3. etc...

And administer needs are:

	1. developers need to know very less about docker to use this.
	2. developers have limited permission when login to docker. 
	3. anyone who does anything should be recorded.

# What is this
Implement Docker CLI on the web.
include:

	1. web shell.  quick exec into a container in k8s.
	2. weblog.  output docker logs on the web page.
	3. attach. attach to a container, used on production.
	4. plugins. help implement user privilege etc.



# Architecture

```
		      --> agent ---> docker --> container		
WebClient ---> proxy |
		      --> agent ---> docker --> container

```

# How to install

    1. run make (make mac_build on mac) to compile.
    2. deploy proxy use deployment , one for a cluster.
    3. deploy agent use daemonset, one for a node.


# Thanks
inspired and copied a lot from projects:
1. https://github.com/aylei/kubectl-debug
2. https://github.com/maoqide/kubeutil 


