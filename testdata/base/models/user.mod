name: User
fields:
  ID:
    type: UUID
    attributes:
      - mandatory
  Email:
    type: String
    attributes:
      - nullable
  Name:
    type: String
identifiers:
  primary: ID
related: {}


