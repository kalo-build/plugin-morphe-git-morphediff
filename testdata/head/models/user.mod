name: User
fields:
  ID:
    type: UUID
    attributes:
      - mandatory
  Email:
    type: String
    attributes:
      - mandatory
  Name:
    type: String
  PhoneNumber:
    type: String
    attributes:
      - nullable
identifiers:
  primary: ID
related:
  ContactInfo:
    type: HasOne

