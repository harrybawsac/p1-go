```markdown
# Data Model: Meter readings and external readings

Entities:

- MeterReading

  - id: uuid / serial (internal)
  - unique_id: string (unique from meter payload)
  - timestamp: timestamptz
  - import_timestamp: timestamptz
  - electricity_delivered: numeric
  - electricity_returned: numeric
  - gas_m3: numeric
  - created_at: timestamptz

- ExternalReading
  - id: uuid / serial
  - meter_reading_id: FK -> MeterReading.id
  - name: text
  - value: text
  - created_at: timestamptz

Relationships:

- MeterReading 1 - N ExternalReading

Validation rules:

- unique_id must be present and globally unique
- timestamp must be parsable as RFC3339 or epoch
- numeric fields must be non-negative when present
```
