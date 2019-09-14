INSERT INTO networks (id, network_name, private_key, port, cidr) 
VALUES (
  '8aeb2ad5-422d-4539-a4e5-5b53f4202c46',
  'TestNet',
  'GNaBxhQICVbdMSTjG5gB/cjthqTVYfuIU1/PQ5PlaUA=',
  '61122',
  '169.254.130.0/24'
);
INSERT INTO peers (id, public_key, peer_name, psk, cidr, network)
VALUES (
  '22eb2ad5-422d-4539-a4e5-5253f1202c46',
  'cNfEWeXuKjazBel/03RkfCgk1yjaX+/V0je5i+3JSF0=',
  'Test',
  'ij89GTql/cqkdBJ86nzvoP6e7KgAbKuXl1Z9ShudG8g=',
  '169.254.130.5/32',
  '8aeb2ad5-422d-4539-a4e5-5b53f4202c46'
);