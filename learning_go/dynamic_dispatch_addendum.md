# Dynamic Dispatch
<!--TODO: Pull info from various sources about dynamic dispatch and why go doesn't have it-->
In-progress
<!-- Dynamic dispatch is a key concept in object-oriented programming, often associated with polymorphism, where the call to an overridden method is resolved at runtime. This contrasts with static dispatch, where the method call is resolved at compile time. -->
<!-- Dynamic Dispatch -->
<!---->
<!-- Dynamic dispatch occurs in languages that support inheritance and method overriding. Here's a basic outline of how it works: -->
<!---->
<!--     Polymorphism: In object-oriented languages, a base class or interface can reference an instance of any derived class. -->
<!---->
<!--     Method Overriding: Derived classes can provide their own implementation for methods of the base class. -->
<!---->
<!--     Runtime Resolution: When a method is called on a polymorphic object, the actual method that gets executed is determined at runtime based on the object's actual class, not the class of the reference variable. This enables one to write more generic and reusable code. -->
<!---->
<!-- Why Go Doesn't Support Dynamic Dispatch -->
<!---->
<!-- Go, while being a modern language with some object-oriented features, doesn't fully support traditional object-oriented concepts like inheritance and class-based polymorphism. Here's why dynamic dispatch is not a concept in Go: -->
<!---->
<!--     No Inheritance: Go doesn't have inheritance. Instead, it has a composition model using "embedding". Without inheritance, there's no concept of base and derived classes, which are fundamental for dynamic dispatch. -->
<!---->
<!--     Interfaces: Go uses interfaces, but they work differently compared to traditional object-oriented languages. In Go, interfaces are satisfied implicitly, meaning that a type satisfies an interface by implementing its methods. There's no explicit declaration of intent to implement a specific interface. -->
<!---->
<!--     Static Typing and Method Resolution: Go is statically typed, and method calls on interface types are resolved at compile time based on the method set of the interface. While it may seem like dynamic behavior, it's more about type compatibility and less about dynamic dispatch as understood in classical object-oriented languages like Java or C++. -->
<!---->
<!--     Design Philosophy: Go's design philosophy emphasizes simplicity and straightforwardness. The language designers intentionally avoided some of the complexities of traditional object-oriented programming, including dynamic dispatch, to keep the language more predictable and easier to understand. -->
<!---->
<!-- In summary, Go's approach to object-oriented programming is unique. It incorporates interfaces and type embedding instead of classical inheritance and polymorphism, aligning with its goals for simplicity and efficiency. This design choice inherently excludes the use of dynamic dispatch as seen in more traditional object-oriented languages. -->
