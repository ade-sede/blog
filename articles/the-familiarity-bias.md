# The familiarity bias

When teams build software, a concept frequently discussed is _simplicity_. The idea that the way we approach problems can be uncomplicated, straightforward.
- _"This code would be much simpler if we used library Y instead"_
- _"We should make this function generic, it will be simpler to reuse"_
- _"Let's get rid of these 5 if statements using pattern X, it's going to make the code simpler"_

These conversations are reoccuring. Why would anyone think these changes make code simpler? Because the resulting code is less verbose? Because the library is popular? Because we can put a name to the pattern and point to the book where it was first introduced?

The actual reason is obvious but often overlooked: because they have seen a similar code before and people's sense of _what is simple_ is overwhelmingly dominated by familiarity.

This bias is something we as human beings cannot escape. Perception is subjective. We do not have the ability to forget what we know and judge situations objectively.

In this article I want to explore how to write about how to identify this bias in ourselves and in others, as well as how to defuse it.
But first, let's talk about what simplicity is.

## What is simplicity in software?

Simplicity spans several dimensions[^1]:
- Something simple does not drown us under its _cognitive load_
- Something simple is _intuitive_ and _self contained_: no specialised knowledge is required to understand it
- Something simple is _predictable_ and _uniform_: once we understand it we are able to predict how it (and things that resemble it) behave in any given scenario
- Simple things _scale_ and are _composable_: they can be assembled together and remain simple. In contrast, when we assembled complex thing the complexity compounds and explodes, overloading our cognitive ability for even the smallest of _n_

I realise this definition is incomplete and abstract so let's try to illustrate it with an example.

### Example: illustrating simplicity

The three following implementations are functionally equivalent[^2] and yield the same result:

#### Implementation 1: combinatorial logic

The first implementation combines all attributes of a `User` to decide which dashboard to display.  
When reading this code, we are forced to track which combinations have been checked. This only gets worse as we add properties to our `User`. We only have 3 properties and I find it hard to carry all the state around in my head[^3].  

```typescript
type User = {
  hasAnalytics: boolean,
  hasReporting: boolean, 
  isTrialUser: boolean,
}

function getDashboardConfig(user: User): DashboardConfig {
  if (user.hasAnalytics && user.hasReporting && !user.isTrialUser) {
    return fullAnalyticsDashboard();
  }

  if (user.hasAnalytics && !user.hasReporting && !user.isTrialUser) {
    return analyticsDashboard();
  }

  if (!user.hasAnalytics && user.hasReporting && !user.isTrialUser) {
    return reportingDashboard();
  }

  if (user.hasAnalytics && user.hasReporting && user.isTrialUser) {
    return trialReportingDashboard();
  }

  if (user.hasAnalytics && !user.hasReporting && user.isTrialUser) {
    return trialAnalyticsDashboard();
  }

  if (!user.hasAnalytics && !user.hasReporting && user.isTrialUser) {
    return basicTrialDashboard();
  }

  return defaultDashboard();
}
```

#### Implementation 2: strategy pattern 

This second implementation has a strong OOP feel. Someone reading `getDashboardConfig` has no idea what the function does. We can accept it as a black box, or we can try to traverse the hierarchy of objects to understand how it works. For every attribute that is added the hierarchy grows longer and reading becomes more strenuous. I just wrote this code 30 seconds ago and I can't read it anymore without being overwhelmed...

```typescript
type User = {
  hasAnalytics: boolean,
  hasReporting: boolean, 
  isTrialUser: boolean
}

interface DashboardStrategy {
  execute(): DashboardConfig;
  withFallback(fallback: DashboardStrategy): DashboardStrategy;
}

class FullAnalyticsDashboardStrategy implements DashboardStrategy {
  execute(): DashboardConfig {
    return fullAnalyticsDashboard();
  }
  
  withFallback(fallback: DashboardStrategy): DashboardStrategy {
    return this;
  }
}

class AnalyticsDashboardStrategy implements DashboardStrategy {
  execute(): DashboardConfig {
    return analyticsDashboard();
  }
  
  withFallback(fallback: DashboardStrategy): DashboardStrategy {
    return this;
  }
}

class ReportingDashboardStrategy implements DashboardStrategy {
  execute(): DashboardConfig {
    return reportingDashboard();
  }
  
  withFallback(fallback: DashboardStrategy): DashboardStrategy {
    return this;
  }
}

class TrialAnalyticsDashboardStrategy implements DashboardStrategy {
  execute(): DashboardConfig {
    return trialAnalyticsDashboard();
  }
  
  withFallback(fallback: DashboardStrategy): DashboardStrategy {
    return this;
  }
}

class DefaultDashboardStrategy implements DashboardStrategy {
  execute(): DashboardConfig {
    return defaultDashboard();
  }
  
  withFallback(fallback: DashboardStrategy): DashboardStrategy {
    return this;
  }
}

class NoOpStrategy implements DashboardStrategy {
  private fallback: DashboardStrategy;
  
  constructor(fallback: DashboardStrategy) {
    this.fallback = fallback;
  }
  
  execute(): DashboardConfig {
    return this.fallback.execute();
  }
  
  withFallback(fallback: DashboardStrategy): DashboardStrategy {
    return fallback;
  }
}

class DashboardStrategyFactory {
  createStrategy(user: User): DashboardStrategy {
    if (user.hasAnalytics && user.hasReporting && !user.isTrialUser) {
      return new FullAnalyticsDashboardStrategy();
    }
    
    if (user.hasAnalytics && !user.hasReporting && !user.isTrialUser) {
      return new AnalyticsDashboardStrategy();
    }
    
    if (!user.hasAnalytics && user.hasReporting && !user.isTrialUser) {
      return new ReportingDashboardStrategy();
    }
    
    if (user.hasAnalytics && !user.hasReporting && user.isTrialUser) {
      return new TrialAnalyticsDashboardStrategy();
    }
    
    return new NoOpStrategy(new DefaultDashboardStrategy());
  }
}

function getDashboardConfig(user: User): DashboardConfig {
  const dashboardStrategy = new DashboardStrategyFactory()
    .createStrategy(user)
    .withFallback(new DefaultDashboardStrategy());
    
  return dashboardStrategy.execute();
}
```
#### Implementation 3: linear logic

In this third implementation we have adapted the way we model our `User`.  
It has enabled us to apply simple, linear logic. Every `archetype` is mapped to one and only one dashboard. There could be an infinity of archetypes, this function would still be as easy to understand[^4].


```javascript
type User = {
  archetype: 'power-user' | 'analyst' | 'reporter' | 
             'trial-user' | 'basic-trial'
}

function getDashboardConfig(user) {
  switch (user.archetype) {
    case 'power-user':
      return fullAnalyticsDashboard();
    case 'analyst':
      return analyticsDashboard();
    case 'reporter':
      return reportingDashboard();
    case 'trial-user':
      return trialAnalyticsDashboard();
    case 'basic-trial':
      return basicTrialDashboard();
    default:
      return defaultDashboard();
  }
}
```

### Why do these exemples read 

## DRAFT NOTES (WIP)

The most common source of complexity in programming is abstraction.




complexity explodes, compounds
m + n or m * n ?
cognitive bounds
knowledge bounds

??
Pushing back on solutions someone perceives as simple can be a complicated conversation. Because in their eyes, you are saying no to simplicity. Questioning  
??

increasing your team's tolerance for complexity?
turning complex into simple?


In this article I want to to try and define simplicity and discuss how to spot & defuse this familiarity bias.
people link abstraction to intelligence. saying no to their abstraction can cause them to feel inadequate


Can feel like a personal attack?

Simplicity = less verbose?
Simplicity = pattern has a name?


debating whether something is simple or not is a fool's errand
it becomes 'personal' -> the guy who doesn't want simple is a bad dev
let's instead shift the conversation to some more concrete attributes of the proposed solutions
blanket word


Does that mean we should get rid of 'patterns' and 'language' and everything that cannot be instantly understood by a toddler?
Admit that there exists situations where simplicity is the what we should optimise for.

the word 'simple' is hiding stuff and reducing quality of conversation / distracting from specific properties of the solution

In the name of simplicity, people push in all kinds of direction (monads? :D)
And saying no is not an option: What kind of incompetent software developer would say not to simplicity?!

We are lacking common definitions.

Everyone thinks what they are doing is simple.
Because it is simple _to them_.

The inability to let go of our own POV is what is blocking us down.

We need forced functions to help us get rid of our own bias.

If its simple to you but not to your teammates -> what knowledge are they missing to turn it into something simple?

What feels 'good' 'bad' -> software instinct? craftsman instinct?

Broadcast what simple is and should be <- Senior IC work?

cognitive ability? knowledge? -> upper bound of simplicity?
expand upper bound by sharing knowledge at the right time

universal simplicity (== regardless of individual?) Does it exist?

Sometimes simplicity is about constraining the complexity and putting it in the right place
The right abstraction will reduce complexity.
The wrong abstraction will increase complexity.
Finding the right abstraction is hard, when we don't know what it is might as well not abstract at all.
You can find an abstraction that feels right _to people who have spent enough time in the problem space_ -> knowledge + familiarity = upper bound


defusing -> 'which of the things you are familiar with are applicable in our case'

you want people to do simple but you don't want them to freeze complexity arises

should we optimise for read or for write ?

[^1]: While writing this article I stumbled upon the excelent talk ["Simple Made Easy"](https://www.youtube.com/watch?v=SxdOUGdseq4) by Rich Hickey. I don't agree with all of the takeaways, but it is definitely the best explanation of what simplicy is I have ever heard. So much so I feel stupid writing this section of the article... Highly recommend watching.
[^2]: If you squint really hard ...
[^3]: Point being: no matter what your personal threshold is it will get exceeded very quickly, because combinatorial logic is explosive
[^4]: The truth is the logic to decide the archetype of a user probably still looks somewhat like implementation 1. We haven't exactly killed the complexity, we just moved it. I swear finding good examples is hard...
