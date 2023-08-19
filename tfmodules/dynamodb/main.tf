resource "aws_dynamodb_table" "shopit_item" {
  name           = "shopit_item"
  hash_key       = "Id"
  read_capacity  = 30
  write_capacity = 30

  attribute {
    name = "Id"
    type = "S"
  }
}

output "shopit_item_table_name" {
  value = aws_dynamodb_table.shopit_item.name

  depends_on = [aws_dynamodb_table.shopit_item]
}
